package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

type HttpMixerOptions struct {
	source      *string
	concurrency *int
	timeout     *int
	redirect    *bool
	testHttp    *bool
	testHttps   *bool
	testTrace   *bool
}

type HttpMixerResult struct {
	statusCode int
	location   string
	url        string
	// trace      bool
}

type HttpMixer struct {
	source  io.Reader
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
	feedChannel := make(chan string)
	outWG := &sync.WaitGroup{}
	feedWG := &sync.WaitGroup{}

	for i := 0; i < *h.options.concurrency; i++ {
		feedWG.Add(1)
		go func() {
			for url := range feedChannel {
				result, err := h.client.request(&url, "GET")
				if err != nil {
					continue
				}

				outChannel <- result
			}
			feedWG.Done()
		}()
	}

	outWG.Add(1)
	go func() {
		defer outWG.Done()
		for o := range outChannel {
			f(o)
		}
	}()

	go func() {
		feedWG.Wait()
		close(outChannel)
	}()

	h.feed(feedChannel)

	close(feedChannel)
	outWG.Wait()
}

func (h *HttpMixer) feed(feedChannel chan string) {
	scanner := bufio.NewScanner(h.source)
	for scanner.Scan() {
		urls := h.urlsWthProtocols(scanner.Text())
		for _, url := range urls {
			feedChannel <- url
		}
	}
}

func (h *HttpMixer) urlsWthProtocols(url string) []string {
	result := []string{}

	if *h.options.testHttp {
		if !strings.HasPrefix(url, "http") {
			result = append(result, "http://"+url)
		} else {
			result = append(result, url)
		}
	}

	if *h.options.testHttps {
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
