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

type selectedCodes []int

type statusFilter struct {
	showAll       bool
	onlyInfo      bool
	onlySuccess   bool
	onlyClientErr bool
	onlyServerErr bool
	selected      selectedCodes
}

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
		options: opts,
		output:  printResult,
		// output:  &output{outputFunction: printResult},
		client: getClient(&opts.noRedirect, &opts.timeout),
	}, nil
}

type Summary map[string]map[string]int

var summaryData = make(Summary)

func (o *HttpMixerOptions) reprSource() string {
	if o.source == "" {
		return Blue("source: ") + Green("stdout")
	}
	return Blue("source: ") + Green("stdout")
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

type resultsFunc func(result *HttpMixerResult)

// Output represents way of processing HttpMixerResult outputs. By design outputs
// are redirected to stdout and formatted by internal function. Formatting function
// can be changed by changeOutputFunction method.
// type output struct {
// 	// toChannel      bool
// 	// outputChannel  chan HttpMixerResult
// 	outputFunction resultsFunc
// }

type HttpMixer struct {
	source  io.ReadCloser
	client  *HttpClient
	options *HttpMixerOptions
	// output  *output
	output resultsFunc
}

// func (h *HttpMixer) setOutputChannel(out chan HttpMixerResult) {
// 	h.options.pipe = true
// 	h.output.toChannel = true
// 	h.output.outputChannel = out
// }

func (h *HttpMixer) Start() {
	start := time.Now()

	if !h.options.pipe {
		printInfo(h.options)
	}

	outChannel := make(chan *HttpMixerResult)
	feedChannel := make(chan *feedData)
	outWG := &sync.WaitGroup{}
	feedWG := &sync.WaitGroup{}

	saveOutput := false
	var outputFile io.WriteCloser
	var outputWriter *bufio.Writer

	if h.options.output != "" {
		saveOutput = true
		outputFile = createFile(h.options.output, 7)
		outputWriter = bufio.NewWriter(outputFile)
	}

	for i := 0; i < h.options.concurrency; i++ {
		feedWG.Add(1)
		go func() {
			for feed := range feedChannel {
				result, err := h.client.request(&feed.url, feed.method)
				if err != nil {
					// fmt.Println(err.Error()) // debug
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
		h.output = func(result *HttpMixerResult) {
			fmt.Println(result.url)
		}
	}

	outWG.Add(1)
	go func() {
		defer outWG.Done()
		for o := range outChannel {
			o.description = resolveCodeDescription(o.statusCode, h.options.statusFilter)
			o.url = strings.TrimRight(o.url, "/")

			if len(h.options.statusFilter.selected) > 0 {
				if !intSliceContains(h.options.statusFilter.selected, o.statusCode) {
					continue
				}
			}

			// if h.output.toChannel {
			// 	// can be set by setOutputChannel()
			// 	h.output.outputChannel <- *o
			// } else {
			// 	// can be changed by changeOutputFunction()
			// 	h.output.outputFunction(o)
			// }

			h.output(o)

			if saveOutput {
				_, err := outputWriter.WriteString(
					fmt.Sprintf("%s,%d,%s\n", o.method, o.statusCode, o.url),
				)

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

// setSource sets slice of strings (s) as source. The previous
// source will be closed.
func (h *HttpMixer) setSource(s []string) {
	h.source.Close()
	h.source = ioutil.NopCloser(strings.NewReader(strings.Join(s, "\n")))
}

// changeOutputFunction sets custom function (f) which will be responsible
// for handling results included in HttpMixerResult.
// func (h *HttpMixer) changeOutputFunction(f resultsFunc) {
// 	h.output.outputFunction = f
// }

// wthProtocols prepares slice of two strings, URLs
// with http and https protocols accordingly.
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
