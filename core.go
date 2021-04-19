package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
		urls := h.urlsWthProtocols(scanner.Text())
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

func (h *HttpMixer) urlsWthProtocols(url string) []string {
	result := []string{}

	if !*h.options.skipHttp {
		if !strings.HasPrefix(url, "http") {
			result = append(result, "http://"+url)
		} else {
			result = append(result, url)
		}
	}

	if !*h.options.skipHttps {
		if !strings.HasPrefix(url, "https") {
			result = append(result, "https://"+url)
		} else {
			result = append(result, url)
		}
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
