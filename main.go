package main

import (
	"flag"
	"fmt"
)

func main() {
	source := flag.String("source", "", "Path to file with URL's to test (default: stdin)")
	concurrency := flag.Int("concurrency", 100, "Concurrency level (default: 100)")
	timeout := flag.Int("timeout", 20, "Timeout in seconds (default: 15)")
	redirect := flag.Bool("redirect", true, "Follow redirections? (default: true)")
	skipHttp := flag.Bool("skip-http", false, "Skip testing HTTP protocol (default: false)")
	skipHttps := flag.Bool("skip-https", false, "Skip testing HTTPS protocol (default: false)")
	testTrace := flag.Bool("test-trace", false, "Test TRACE method? (default: false)")

	flag.Parse()

	o := &HttpMixerOptions{
		source:      source,
		concurrency: concurrency,
		redirect:    redirect,
		timeout:     timeout,
		skipHttp:    skipHttp,
		skipHttps:   skipHttps,
		testTrace:   testTrace,
	}

	f := func(o *HttpMixerResult) {
		fmt.Println(o)
	}

	mixer := NewHttpMixer(o)
	mixer.Start(f)
}
