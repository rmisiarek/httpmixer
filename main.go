package main

import (
	"flag"
	"fmt"
)

func main() {
	source := flag.String("source", "", "Path to file with URL's to test (default: stdin)")
	concurrency := flag.Int("concurrency", 100, "Concurrency level (default: 100)")
	timeout := flag.Int("timeout", 15, "Timeout in seconds (default: 15)")
	redirect := flag.Bool("redirect", true, "Follow redirections? (default: true)")
	testHttp := flag.Bool("test-http", true, "Test HTTP protocol? (default: true)")
	testHttps := flag.Bool("test-https", true, "Test HTTPS protocol? (default: true)")
	testTrace := flag.Bool("test-trace", false, "Test TRACE method? (default: false)")

	flag.Parse()

	o := &HttpMixerOptions{
		source:      source,
		concurrency: concurrency,
		redirect:    redirect,
		timeout:     timeout,
		testHttp:    testHttp,
		testHttps:   testHttps,
		testTrace:   testTrace,
	}

	f := func(o *HttpMixerResult) {
		fmt.Println(o)
	}

	mixer := NewHttpMixer(o)
	mixer.Start(f)
}
