package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpMixerStdinSource(t *testing.T) {
	opts := mixerOptions()
	mixer := NewHttpMixer(&opts)

	// Test whether the function has been set correctly
	f1 := noRedirect(nil, nil)
	f2 := mixer.client.client.CheckRedirect(nil, nil)
	assert.Equal(t, f1, f2)

	// Test whether has been set correctly
	assert.Equal(t, time.Duration(5)*time.Second, mixer.client.client.Timeout)

	// Test access to HttpMixerOptions values
	assert.Equal(t, "", *mixer.options.source)
	assert.Equal(t, 2, *mixer.options.concurrency)
	assert.Equal(t, 5, *mixer.options.timeout)
	assert.Equal(t, false, *mixer.options.redirect)
	assert.Equal(t, false, *mixer.options.skipHttp)
	assert.Equal(t, false, *mixer.options.skipHttps)
	assert.Equal(t, true, *mixer.options.testTrace)
}

func TestNewHttpMixerFileSource(t *testing.T) {
	// Create temporary source file
	tmpSourceFile := createTemporarySourceFile()
	defer os.Remove(tmpSourceFile.Name())

	// Set temporary file to HttpMixerOptions
	opts := mixerOptions()
	source := tmpSourceFile.Name()
	opts.source = &source

	// Success scenario - read from file
	mixer := NewHttpMixer(&opts)
	scanner := bufio.NewScanner(mixer.source)
	for scanner.Scan() {
		assert.Equal(t, "www.google.com", scanner.Text())
	}

	// Failure scenario - there is no such file, should panic
	source = "source-that-not-exist.txt"
	opts.source = &source
	assert.Panics(t, func() { NewHttpMixer(&opts) })
}

func createTemporarySourceFile() *os.File {
	tmpFile, err := ioutil.TempFile(".", "tmp-*-source.txt")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = tmpFile.Write([]byte("www.google.com")); err != nil {
		log.Fatal(err)
	}

	return tmpFile
}

func TestUrlsWthProtocols(t *testing.T) {
	_true := true
	_false := false

	opts := mixerOptions()
	mixer := NewHttpMixer(&opts)

	expected := []string{"http://www.example.com", "https://www.example.com"}

	// Three tests for both http and https
	result := mixer.wthProtocols("www.example.com")
	assert.Equal(t, expected, result)

	result = mixer.wthProtocols("http://www.example.com")
	assert.Equal(t, expected, result)

	result = mixer.wthProtocols("https://www.example.com")
	assert.Equal(t, expected, result)

	// Three tests for https only
	mixer.options.skipHttp = &_true

	result = mixer.wthProtocols("http://www.example.com")
	assert.Equal(t, []string{"https://www.example.com"}, result)

	result = mixer.wthProtocols("https://www.example.com")
	assert.Equal(t, []string{"https://www.example.com"}, result)

	result = mixer.wthProtocols("www.example.com")
	assert.Equal(t, []string{"https://www.example.com"}, result)

	// Three tests for http only
	mixer.options.skipHttp = &_false
	mixer.options.skipHttps = &_true

	result = mixer.wthProtocols("http://www.example.com")
	assert.Equal(t, []string{"http://www.example.com"}, result)

	result = mixer.wthProtocols("https://www.example.com")
	assert.Equal(t, []string{"http://www.example.com"}, result)

	result = mixer.wthProtocols("www.example.com")
	assert.Equal(t, []string{"http://www.example.com"}, result)
}

func mixerOptions() HttpMixerOptions {
	source := ""
	concurrency := 2
	timeout := 5
	redirect := false
	skipHttp := false
	skipHttps := false
	testTrace := true
	showAll := true
	onlyInfo := false
	onlySuccess := false
	onlyClientErr := false
	onlyServerErr := false

	filter := statusFilter{
		showAll:       &showAll,
		onlyInfo:      &onlyInfo,
		onlySuccess:   &onlySuccess,
		onlyClientErr: &onlyClientErr,
		onlyServerErr: &onlyServerErr,
	}

	opts := HttpMixerOptions{
		source:       &source,
		concurrency:  &concurrency,
		timeout:      &timeout,
		redirect:     &redirect,
		skipHttp:     &skipHttp,
		skipHttps:    &skipHttps,
		testTrace:    &testTrace,
		statusFilter: &filter,
	}

	return opts
}
