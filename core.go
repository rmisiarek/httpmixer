package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var c = &http.Client{
	Timeout:       15 * time.Second,
	CheckRedirect: noRedirect,
}

type HttpMixerOpts struct {
	source    *string
	redirect  *bool
	testHttp  bool
	testHttps bool
}

type HttpMixerRes struct {
	statusCode int
	location   string
	url        string
	trace      bool
}

type HttpMixer struct {
	source io.Reader
	client *http.Client
}

func NewHttpMixer(opts *HttpMixerOpts) *HttpMixer {
	return &HttpMixer{
		source: openStdinOrFile(opts.source),
		client: getClient(opts.redirect),
	}
}

type resultF func(result *HttpMixerRes)

func (h *HttpMixer) Start(f resultF) {
	outChannel := make(chan *HttpMixerRes)
	feedChannel := make(chan string)
	outWG := &sync.WaitGroup{}
	feedWG := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		feedWG.Add(1)
		go func() {
			for url := range feedChannel {
				u := url

				resp, code, err := checkStatus(u)
				if err != nil {
					continue
				}

				outChannel <- &HttpMixerRes{location: resp, statusCode: code}
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
		urls := prepareUrlsWithSchema(scanner.Text())
		for _, url := range urls {
			feedChannel <- url
		}
	}
}

func (h *HttpMixer) request(url *string) (*HttpMixerRes, error) {
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	location := ""
	// _location, err := resp.Location()
	// if err != nil {
	// 	location = ""
	// } else {
	// 	location = _location.String()
	// }

	result := &HttpMixerRes{
		statusCode: resp.StatusCode,
		location:   location,
		url:        resp.Request.URL.String(),
	}

	return result, nil
}

// func (h *HttpMixer) urlsWthProtocols(url *string) []string {
// 	result := []string{}

// 	if h.opts.testHttp {
// 		if !strings.HasPrefix(*url, "http") {
// 			result = append(result, "http://"+*url)
// 		} else {
// 			result = append(result, *url)
// 		}
// 	}

// 	if h.opts.testHttps {
// 		if !strings.HasPrefix(*url, "https") {
// 			result = append(result, "https://"+*url)
// 		} else {
// 			result = append(result, *url)
// 		}
// 	}

// 	return result
// }

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
