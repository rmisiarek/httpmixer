package main

import (
	"flag"
)

func main() {
	source := flag.String("source", "", "Path to file with URL's to test (default: stdin)")
	output := flag.String("output", "", "Path to output file (default: stdout)")
	concurrency := flag.Int("concurrency", 100, "Concurrency level (defaqqqqult: 100)")
	timeout := flag.Int("timeout", 30, "Timeout in seconds (default: 30s)")
	noRedirect := flag.Bool("no-redirect", false, "Don't follow redirections (default: false)")
	skipHttp := flag.Bool("skip-http", false, "Skip testing HTTP protocol (default: false)")
	skipHttps := flag.Bool("skip-https", false, "Skip testing HTTPS protocol (default: false)")
	testTrace := flag.Bool("test-trace", false, "Test TRACE method? (default: false)")
	onlyInfo := flag.Bool("info", false, "Filter only informational statuses (default: false)")
	onlySuccess := flag.Bool("success", false, "Filter only success statuses (default: false)")
	onlyClientErr := flag.Bool("client-error", false, "Filter only client error statuses (default: false)")
	onlyServerErr := flag.Bool("server-error", false, "Filter only server error statuses (default: false)")

	flag.Parse()

	options := &HttpMixerOptions{
		source:      source,
		output:      output,
		concurrency: concurrency,
		noRedirect:  noRedirect,
		timeout:     timeout,
		skipHttp:    skipHttp,
		skipHttps:   skipHttps,
		testTrace:   testTrace,
		statusFilter: &statusFilter{
			onlyInfo:      onlyInfo,
			onlySuccess:   onlySuccess,
			onlyClientErr: onlyClientErr,
			onlyServerErr: onlyServerErr,
		},
	}

	// f := func(o *HttpMixerResult) {
	// 	fmt.Println(o)
	// }

	mixer := NewHttpMixer(options)
	printInfo(options)

	mixer.Start(printResult)
}
