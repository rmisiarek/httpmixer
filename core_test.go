package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpMixerStdinSource(t *testing.T) {
	opts := mixerOptions()
	opts.noRedirect = true

	mixer, _ := NewHttpMixer(&opts)

	// Test whether the function has been set correctly
	f1 := noRedirectF(nil, nil)
	f2 := mixer.client.client.CheckRedirect(nil, nil)
	assert.Equal(t, f1, f2)

	// Test whether timeout has been set correctly
	assert.Equal(t, time.Duration(5)*time.Second, mixer.client.client.Timeout)

	// Test access to HttpMixerOptions values
	assert.Equal(t, "", mixer.options.source)
	assert.Equal(t, 2, mixer.options.concurrency)
	assert.Equal(t, 5, mixer.options.timeout)
	assert.Equal(t, true, mixer.options.noRedirect)
	assert.Equal(t, false, mixer.options.skipHttp)
	assert.Equal(t, false, mixer.options.skipHttps)
	assert.Equal(t, true, mixer.options.testTrace)
}

func TestNewHttpMixerFileSource(t *testing.T) {
	// Create temporary source file
	tmpSourceFile := createTemporarySourceFile()
	defer os.Remove(tmpSourceFile.Name())

	// Set temporary file to HttpMixerOptions
	opts := mixerOptions()
	source := tmpSourceFile.Name()
	opts.source = source

	// Success scenario - read from file
	mixer, _ := NewHttpMixer(&opts)
	scanner := bufio.NewScanner(mixer.source)
	for scanner.Scan() {
		assert.Equal(t, "www.google.com", scanner.Text())
	}

	// Failure scenario - there is no such file, should return err
	source = "source-that-not-exist.txt"
	opts.source = source
	mixer, err := NewHttpMixer(&opts)
	assert.Nil(t, mixer)
	assert.NotNil(t, err)
}

func TestNewHttpMixerStatusFilter(t *testing.T) {
	opts := mixerOptions()
	m, _ := NewHttpMixer(&opts)
	assert.Equal(t, true, m.options.statusFilter.showAll)

	// When any of the filters is set to true, then showAll
	// should be set to false, just to show only filtered results
	opts.statusFilter.onlySuccess = true
	m, _ = NewHttpMixer(&opts)
	assert.Equal(t, false, m.options.statusFilter.showAll)
}

func TestUrlsWthProtocols(t *testing.T) {
	opts := mixerOptions()
	mixer, _ := NewHttpMixer(&opts)

	expected := []string{"http://www.example.com", "https://www.example.com"}

	// Three tests for both http and https
	result := mixer.wthProtocols("www.example.com")
	assert.Equal(t, expected, result)

	result = mixer.wthProtocols("http://www.example.com")
	assert.Equal(t, expected, result)

	result = mixer.wthProtocols("https://www.example.com")
	assert.Equal(t, expected, result)

	// Three tests for https only
	mixer.options.skipHttp = true

	result = mixer.wthProtocols("http://www.example.com")
	assert.Equal(t, []string{"https://www.example.com"}, result)

	result = mixer.wthProtocols("https://www.example.com")
	assert.Equal(t, []string{"https://www.example.com"}, result)

	result = mixer.wthProtocols("www.example.com")
	assert.Equal(t, []string{"https://www.example.com"}, result)

	// Three tests for http only
	mixer.options.skipHttp = false
	mixer.options.skipHttps = true

	result = mixer.wthProtocols("http://www.example.com")
	assert.Equal(t, []string{"http://www.example.com"}, result)

	result = mixer.wthProtocols("https://www.example.com")
	assert.Equal(t, []string{"http://www.example.com"}, result)

	result = mixer.wthProtocols("www.example.com")
	assert.Equal(t, []string{"http://www.example.com"}, result)
}

func TestHttpMixerOptionsRepr(t *testing.T) {
	o := mixerOptions()

	// assert.Equal(t, Blue("source: ")+Green("stdin"), o.reprSource())
	assert.Equal(t, Blue("output: ")+Green("stdout"), o.reprOutput())
	assert.Equal(t, Blue("concurrency: ")+Green(strconv.Itoa(2)), o.reprConcurenncy())
	assert.Equal(t, Blue("timeout: ")+Green(strconv.Itoa(5)), o.reprTimeout())
	assert.Equal(t, Blue("pipe: ")+Green("on"), o.reprPipe())
	assert.Equal(t, Blue("redirect: ")+Green("on"), o.reprRedirect())
	assert.Equal(t, Blue("HTTP: ")+Green("on"), o.reprSkipHttp())
	assert.Equal(t, Blue("HTTPS: ")+Green("on"), o.reprSkipHttps())
	assert.Equal(t, Blue("trace: ")+Green("on"), o.reprTestTrace())
	assert.Equal(t, Blue("filter: ")+Green("all"), o.reprStatusFilter())

	source := "/tmp/file.txt"
	output := "/tmp/results.txt"

	o.source = source
	o.output = output
	o.pipe = true
	o.noRedirect = true
	o.skipHttp = true
	o.skipHttps = true
	o.testTrace = false
	o.statusFilter.showAll = false
	o.statusFilter.onlyClientErr = true
	o.statusFilter.onlyInfo = true
	o.statusFilter.onlyServerErr = true
	o.statusFilter.onlySuccess = true

	// assert.Equal(t, Blue("source: ")+Green("/tmp/file.txt"), o.reprSource())
	assert.Equal(t, Blue("output: ")+Green("/tmp/results.txt"), o.reprOutput())
	assert.Equal(t, Blue("pipe: ")+Green("on"), o.reprPipe())
	assert.Equal(t, Blue("redirect: ")+Red("off"), o.reprRedirect())
	assert.Equal(t, Blue("HTTP: ")+Red("off"), o.reprSkipHttp())
	assert.Equal(t, Blue("HTTPS: ")+Red("off"), o.reprSkipHttps())
	assert.Equal(t, Blue("trace: ")+Red("off"), o.reprTestTrace())
	assert.Equal(t, Blue("filter: ")+Green("info, success, client error, server error"), o.reprStatusFilter())
}

func TestOpenStdinOrFile(t *testing.T) {
	reader, err := openStdinOrFile("no-such-thing")
	assert.Nil(t, reader)
	assert.NotNil(t, err)

	file, err := ioutil.TempFile(".", "tmp-*-file.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.Remove(file.Name())

	reader, err = openStdinOrFile(file.Name())
	assert.NotNil(t, reader)
	assert.Nil(t, err)
}

func TestSourceFromSlice(t *testing.T) {
	source := []string{"source1", "source2"}

	options := &HttpMixerOptions{
		pipe:         true,
		statusFilter: &statusFilter{},
	}

	mixer, err := NewHttpMixer(options)
	assert.Nil(t, err)

	mixer.setSource(source)

	scanner := bufio.NewScanner(mixer.source)
	for scanner.Scan() {
		assert.Equal(t, true, stringSliceContains(source, scanner.Text()))
	}
}

func TestCreateFile(t *testing.T) {
	// Catch loggers output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	f := createFile("/tmp/results-file.txt", 0)
	assert.Equal(t, "/tmp/results-file.txt", f.Name())

	logExpected := Red(fmt.Sprintf(
		"%s exists and will be overwritten. Are you sure? %d seconds to GO\n",
		f.Name(), 0,
	))

	assert.Contains(t, buf.String(), logExpected)

	assert.Panics(t, func() { createFile("tmp/there-is-no-results-file.txt", 0) })
}

func TestFileExists(t *testing.T) {
	path := "/tmp/results.txt"
	exist := fileExists(path)
	assert.Equal(t, false, exist)

	file, _ := ioutil.TempFile("/tmp", "*.txt")
	exist = fileExists(file.Name())
	assert.Equal(t, true, exist)
}

func createTemporarySourceFile() *os.File {
	file, err := ioutil.TempFile(".", "tmp-*-source.txt")
	if err != nil {
		log.Fatal(err)
	}

	if _, err = file.Write([]byte("www.google.com")); err != nil {
		log.Fatal(err)
	}

	return file
}

func mixerOptions() HttpMixerOptions {
	opts := HttpMixerOptions{
		source:      "",
		output:      "",
		concurrency: 2,
		timeout:     5,
		pipe:        true,
		noRedirect:  false,
		skipHttp:    false,
		skipHttps:   false,
		testTrace:   true,
		statusFilter: &statusFilter{
			showAll:       true,
			onlyInfo:      false,
			onlySuccess:   false,
			onlyClientErr: false,
			onlyServerErr: false,
		},
	}

	return opts
}
