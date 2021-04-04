package main

import (
	"fmt"
)

func main() {
	source := "debug/google.txt"
	redirect := true

	o := &HttpMixerOptions{
		source:    &source,
		redirect:  &redirect,
		testHttp:  true,
		testHttps: true,
	}

	f := func(o *HttpMixerResult) {
		fmt.Println(o)
	}

	mixer := NewHttpMixer(o)
	mixer.Start(f)
}
