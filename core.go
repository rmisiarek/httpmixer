package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type feedData struct {
	method string
	url    string
}

type statusFilter struct {
	showAll       *bool
	onlyInfo      *bool
	onlySuccess   *bool
	onlyClientErr *bool
	onlyServerErr *bool
}

type HttpMixerOptions struct {
	source       *string
	concurrency  *int
	timeout      *int
	redirect     *bool
	skipHttp     *bool
	skipHttps    *bool
	testTrace    *bool
	statusFilter *statusFilter
}

func (o *HttpMixerOptions) reprSource() string {
	if *o.source == "" {
		return Blue("source: ") + Green("stdin")
	}
	return Blue("source: ") + Green(*o.source)
}

func (o *HttpMixerOptions) reprConcurenncy() string {
	return Blue("concurrency: ") + Green(strconv.Itoa(*o.concurrency))
}

func (o *HttpMixerOptions) reprTimeout() string {
	return Blue("timeout: ") + Green(strconv.Itoa(*o.timeout))
}

func (o *HttpMixerOptions) reprRedirect() string {
	if *o.redirect {
		return Blue("redirect: ") + Green("on")
	} else {
		return Blue("redirect: ") + Red("off")
	}
}

func (o *HttpMixerOptions) reprSkipHttp() string {
	if *o.skipHttp {
		return Blue("HTTP: ") + Red("off")
	} else {
		return Blue("HTTP: ") + Green("on")
	}
}

func (o *HttpMixerOptions) reprSkipHttps() string {
	if *o.skipHttps {
		return Blue("HTTPS: ") + Red("off")
	} else {
		return Blue("HTTPS: ") + Green("on")
	}
}

func (o *HttpMixerOptions) reprTestTrace() string {
	if *o.testTrace {
		return Blue("trace: ") + Green("on")
	} else {
		return Blue("trace: ") + Red("off")
	}
}

func (o *HttpMixerOptions) reprStatusFilter() string {
	result := []string{}

	if *o.statusFilter.showAll {
		return Blue("filter: ") + Green("all")
	}
	if *o.statusFilter.onlyInfo {
		result = append(result, "info")
	}
	if *o.statusFilter.onlySuccess {
		result = append(result, "success")
	}
	if *o.statusFilter.onlyClientErr {
		result = append(result, "client error")
	}
	if *o.statusFilter.onlyServerErr {
		result = append(result, "server error")
	}

	return Blue("filter: ") + Green(strings.Join(result, ", "))
}

type HttpMixerResult struct {
	statusCode int
	url        string
	method     string
	location   string
}

type HttpMixer struct {
	source  io.ReadCloser
	client  *HttpClient
	options *HttpMixerOptions
}

func NewHttpMixer(opts *HttpMixerOptions) *HttpMixer {
	return &HttpMixer{
		source:  openStdinOrFile(opts.source),
		client:  getClient(opts.redirect, opts.timeout),
		options: opts,
	}
}

type resultF func(result *HttpMixerResult)

func (h *HttpMixer) Start(f resultF) {
	outChannel := make(chan *HttpMixerResult)
	feedChannel := make(chan *feedData)
	outWG := &sync.WaitGroup{}
	feedWG := &sync.WaitGroup{}

	for i := 0; i < *h.options.concurrency; i++ {
		feedWG.Add(1)
		go func() {
			for feed := range feedChannel {
				result, err := h.client.request(&feed.url, feed.method)
				if err != nil {
					fmt.Println(err.Error())
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

	outWG.Add(1)
	go func() {
		defer outWG.Done()
		for o := range outChannel {
			_, found := resolveCategory(o.statusCode, h.options.statusFilter)
			if found {
				// fmt.Println(category, o.statusCode)
				f(o)
			}
		}
	}()

	h.feed(feedChannel)

	close(feedChannel)
	outWG.Wait()
}

func (h *HttpMixer) feed(feedChannel chan *feedData) {
	scanner := bufio.NewScanner(h.source)
	for scanner.Scan() {
		urls := h.wthProtocols(scanner.Text())
		for _, url := range urls {
			feedChannel <- &feedData{
				method: "GET",
				url:    url,
			}
			if *h.options.testTrace {
				feedChannel <- &feedData{
					method: "TRACE",
					url:    url,
				}
			}
		}
	}
}

func (h *HttpMixer) wthProtocols(url string) []string {
	result := []string{}

	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	if !*h.options.skipHttp {
		result = append(result, "http://"+url)
	}

	if !*h.options.skipHttps {
		result = append(result, "https://"+url)
	}

	return result
}

func openStdinOrFile(inputs *string) io.ReadCloser {
	r := os.Stdin

	if *inputs != "" {
		r = openFile(*inputs)
	}

	return r
}

func openFile(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	return file
}
