package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type feedData struct {
	method string
	url    string
}

type statusFilter struct {
	showAll       bool
	onlyInfo      bool
	onlySuccess   bool
	onlyClientErr bool
	onlyServerErr bool
	// onlySelected  bool
	selected selectedCodes
}

type selectedCodes []int

func (sc *selectedCodes) String() string {
	return fmt.Sprintln(*sc)
}

func (sc *selectedCodes) Set(value string) error {
	for _, code := range strings.Split(value, ",") {
		c, err := strconv.Atoi(code)
		if err != nil {
			continue
		}
		*sc = append(*sc, c)
	}

	return nil
}

type HttpMixerOptions struct {
	source       string
	output       string
	concurrency  int
	timeout      int
	pipe         bool
	noRedirect   bool
	skipHttp     bool
	skipHttps    bool
	testTrace    bool
	statusFilter *statusFilter
}

type Summary map[string]map[string]int

var summaryData = make(Summary)

func (o *HttpMixerOptions) reprSource() string {
	if o.source == "" {
		return Blue("source: ") + Green("stdin")
	}
	return Blue("source: ") + Green(o.source)
}

func (o *HttpMixerOptions) reprOutput() string {
	if o.output == "" {
		return Blue("output: ") + Green("stdout")
	}
	return Blue("output: ") + Green(o.output)
}

func (o *HttpMixerOptions) reprConcurenncy() string {
	return Blue("concurrency: ") + Green(strconv.Itoa(o.concurrency))
}

func (o *HttpMixerOptions) reprTimeout() string {
	return Blue("timeout: ") + Green(strconv.Itoa(o.timeout))
}

func (o *HttpMixerOptions) reprPipe() string {
	if o.pipe {
		return Blue("pipe: ") + Green("on")
	} else {
		return Blue("pipe: ") + Red("off")
	}
}

func (o *HttpMixerOptions) reprRedirect() string {
	if !o.noRedirect {
		return Blue("redirect: ") + Green("on")
	} else {
		return Blue("redirect: ") + Red("off")
	}
}

func (o *HttpMixerOptions) reprSkipHttp() string {
	if o.skipHttp {
		return Blue("HTTP: ") + Red("off")
	} else {
		return Blue("HTTP: ") + Green("on")
	}
}

func (o *HttpMixerOptions) reprSkipHttps() string {
	if o.skipHttps {
		return Blue("HTTPS: ") + Red("off")
	} else {
		return Blue("HTTPS: ") + Green("on")
	}
}

func (o *HttpMixerOptions) reprTestTrace() string {
	if o.testTrace {
		return Blue("trace: ") + Green("on")
	} else {
		return Blue("trace: ") + Red("off")
	}
}

func (o *HttpMixerOptions) reprStatusFilter() string {
	result := []string{}

	if o.statusFilter.showAll {
		result = append(result, "all")
	}
	if o.statusFilter.onlyInfo {
		result = append(result, "info")
	}
	if o.statusFilter.onlySuccess {
		result = append(result, "success")
	}
	if o.statusFilter.onlyClientErr {
		result = append(result, "client error")
	}
	if o.statusFilter.onlyServerErr {
		result = append(result, "server error")
	}
	if len(o.statusFilter.selected) > 0 {
		s := fmt.Sprintf("selected: %s", intArrayToString(o.statusFilter.selected))
		result = append(result, s)
	}

	return Blue("filter: ") + Green(strings.Join(result, ", "))
}

type HttpMixerResult struct {
	statusCode  int
	url         string
	method      string
	location    string
	description string
}

type HttpMixer struct {
	source  io.ReadCloser
	client  *HttpClient
	options *HttpMixerOptions
}

func NewHttpMixer(opts *HttpMixerOptions) (*HttpMixer, error) {
	if opts.source != "" && !fileExists(opts.source) {
		return nil, fmt.Errorf("%s does not exist", opts.source)
	}

	source, err := openStdinOrFile(opts.source)
	if err != nil {
		return nil, fmt.Errorf("error while opening %s", opts.source)
	}

	opts.statusFilter.showAll = true
	if opts.statusFilter.onlyInfo ||
		opts.statusFilter.onlySuccess ||
		opts.statusFilter.onlyClientErr ||
		opts.statusFilter.onlyServerErr {
		opts.statusFilter.showAll = false
	}

	return &HttpMixer{
		source:  source,
		client:  getClient(&opts.noRedirect, &opts.timeout),
		options: opts,
	}, nil
}

type resultF func(result *HttpMixerResult)

func (h *HttpMixer) Start(f resultF) {
	start := time.Now()

	outChannel := make(chan *HttpMixerResult)
	feedChannel := make(chan *feedData)
	outWG := &sync.WaitGroup{}
	feedWG := &sync.WaitGroup{}

	saveOutput := false
	var outputFile io.WriteCloser
	var outputWriter *bufio.Writer

	if h.options.output != "" {
		saveOutput = true
		outputFile = createFile(h.options.output, 5)
		outputWriter = bufio.NewWriter(outputFile)
	}

	for i := 0; i < h.options.concurrency; i++ {
		feedWG.Add(1)
		go func() {
			for feed := range feedChannel {
				result, err := h.client.request(&feed.url, feed.method)
				if err != nil {
					// fmt.Println(err.Error())
					continue
				}

				result.method = feed.method
				outChannel <- result
			}

			feedWG.Done()
		}()
	}

	go func() {
		feedWG.Wait()
		close(outChannel)
	}()

	if h.options.pipe {
		f = func(result *HttpMixerResult) {
			fmt.Println(result.url)
		}
	}

	outWG.Add(1)
	go func() {
		defer outWG.Done()
		for o := range outChannel {
			o.description = resolveCodeDescription(o.statusCode, h.options.statusFilter)

			if len(h.options.statusFilter.selected) > 0 {
				if !intSliceContains(h.options.statusFilter.selected, o.statusCode) {
					continue
				}
			}

			f(o)

			if saveOutput {
				_, err := outputWriter.WriteString(o.url + "\n")
				if err != nil {
					log.Fatalf("error while writing to a file: %s", err.Error())
				}
			}

			if !h.options.pipe {
				aggregateSummary(o, h.options.statusFilter.showAll)
			}
		}
	}()

	h.feed(feedChannel)

	close(feedChannel)
	outWG.Wait()

	if saveOutput {
		outputWriter.Flush()
		outputFile.Close()
	}

	took := time.Since(start).Truncate(time.Second)

	if !h.options.pipe {
		printSummary(summaryData, took)
	}
}

func (h *HttpMixer) feed(feedChannel chan *feedData) {
	defer h.source.Close()

	scanner := bufio.NewScanner(h.source)
	for scanner.Scan() {
		urls := h.wthProtocols(scanner.Text())
		for _, url := range urls {
			feedChannel <- &feedData{
				method: "GET",
				url:    url,
			}
			if h.options.testTrace {
				feedChannel <- &feedData{
					method: "TRACE",
					url:    url,
				}
			}
		}
	}
}

func (h *HttpMixer) sourceFromSlice(s []string) {
	h.source = ioutil.NopCloser(strings.NewReader(strings.Join(s, "\n")))
}

func (h *HttpMixer) wthProtocols(url string) []string {
	result := []string{}

	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	if !h.options.skipHttp {
		result = append(result, "http://"+url)
	}

	if !h.options.skipHttps {
		result = append(result, "https://"+url)
	}

	return result
}
