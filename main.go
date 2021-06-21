package main

import (
	"flag"
	"log"
)

func main() {
	options := &HttpMixerOptions{statusFilter: &statusFilter{}}
	flag.StringVar(&options.source, "source", "", "Path to file with URL's to test (default: stdin)")
	flag.StringVar(&options.output, "output", "", "Path to output file (default: stdout)")
	flag.IntVar(&options.concurrency, "concurrency", 100, "Concurrency level (default: 100)")
	flag.IntVar(&options.timeout, "timeout", 30, "Timeout in seconds (default: 30s)")
	flag.BoolVar(&options.pipe, "pipe", false, "Show only filtered URL's (default: false)")
	flag.BoolVar(&options.noRedirect, "no-redirect", false, "Don't follow redirections (default: false)")
	flag.BoolVar(&options.skipHttp, "skip-http", false, "Skip testing HTTP protocol (default: false)")
	flag.BoolVar(&options.skipHttps, "skip-https", false, "Skip testing HTTPS protocol (default: false)")
	flag.BoolVar(&options.testTrace, "test-trace", false, "Test TRACE method? (default: false)")
	flag.BoolVar(&options.statusFilter.onlyInfo, "info", false, "Filter only informational statuses (default: false)")
	flag.BoolVar(&options.statusFilter.onlySuccess, "success", false, "Filter only success statuses (default: false)")
	flag.BoolVar(&options.statusFilter.onlyClientErr, "client-error", false, "Filter only client error statuses (default: false)")
	flag.BoolVar(&options.statusFilter.onlyServerErr, "server-error", false, "Filter only server error statuses (default: false)")
	flag.Var(&options.statusFilter.selected, "select", "Filter only selected statuses, comma-separated list (default: nil)")
	flag.Parse()

	mixer, err := NewHttpMixer(options)
	if err != nil {
		log.Fatalln(Red(err.Error()))
	}

	// src := []string{
	// 	"google.com", "golang.org",
	// }
	// mixer.sourceFromSlice(src)

	printInfo(options)
	mixer.Start(printResult)
}
