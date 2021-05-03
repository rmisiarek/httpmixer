package main

import (
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

var noRedirectF = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func getClient(withoutRedirect *bool, timeout *int) *HttpClient {
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	if *withoutRedirect {
		client.CheckRedirect = noRedirectF
	}

	return &HttpClient{client: client}
}

func (h *HttpClient) request(url *string, method string) (*HttpMixerResult, error) {
	req, err := http.NewRequest(method, *url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	location := ""
	_location, err := resp.Location()
	if err != nil {
		location = ""
	} else {
		location = _location.String()
	}

	result := &HttpMixerResult{
		statusCode: resp.StatusCode,
		location:   location,
		url:        resp.Request.URL.String(),
	}

	return result, nil
}
