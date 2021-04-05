package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoRedirect(t *testing.T) {
	f := noRedirect(nil, nil)
	assert.Equal(t, f, http.ErrUseLastResponse)
}

func TestGetClient(t *testing.T) {
	redirect := false
	timeout := 10

	c := getClient(&redirect, &timeout)

	assert.Equal(t, time.Duration(timeout)*time.Second, c.client.Timeout)
}

func TestRequestWithoutRedirect(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))
	defer s.Close()

	redirect := false
	timeout := 10

	c := getClient(&redirect, &timeout)
	url := s.URL

	resp, err := c.request(&url, "GET")

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.statusCode)
	assert.Equal(t, "", resp.location)
	assert.Equal(t, url, resp.url)
}

func TestRequestWithRedirect(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "new-location")
		w.WriteHeader(http.StatusMovedPermanently)
		fmt.Fprintln(w, "")
	}))
	defer s.Close()

	redirect := false
	timeout := 10

	c := getClient(&redirect, &timeout)
	url := s.URL

	resp, err := c.request(&url, "GET")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMovedPermanently, resp.statusCode)
	assert.Equal(t, url+"/new-location", resp.location)
	assert.Equal(t, url, resp.url)
}

func TestRequestFailure(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))
	defer s.Close()

	redirect := false
	timeout := 10
	url := ""

	c := getClient(&redirect, &timeout)

	// Failure - unsupported protocol scheme ""
	_, err := c.request(&url, "GET")
	assert.NotNil(t, err)

	// Failure - invalid method
	_, err = c.request(&url, ";")
	assert.NotNil(t, err)
}
