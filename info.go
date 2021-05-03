package main

import (
	"fmt"
)

const title = `
 _      _    _                      _ 
| |    | |  | |                    (_)
| |__  | |_ | |_  _ __   _ __ ___   _ __  __ ___  _ __ 
| '_ \ | __|| __|| '_ \ | '_ ' _ \ | |\ \/ // _ \| '__|
| | | || |_ | |_ | |_) || | | | | || | >  <|  __/| |
|_| |_| \__| \__|| .__/ |_| |_| |_||_|/_/\_\\___||_|
                 | |
                 |_|
`

func printInfo(o *HttpMixerOptions) {
	fmt.Printf("%s\n", Cyan(title))
	fmt.Printf(">> %s\n", o.reprSource())
	fmt.Printf(">> %s\n", o.reprOutput())
	fmt.Printf(">> %s\n", o.reprSkipHttps())
	fmt.Printf(">> %s\n", o.reprSkipHttp())
	fmt.Printf(">> %s\n", o.reprTestTrace())
	fmt.Printf(">> %s\n", o.reprRedirect())
	fmt.Printf(">> %s\n", o.reprTimeout())
	fmt.Printf(">> %s\n", o.reprConcurenncy())
	fmt.Printf(">> %s\n\n", o.reprStatusFilter())
	fmt.Printf("%s\t\t%s\t%s\t%s%s\n", "METHOD", "STATUS CODE", "DESCRIPTION", "\t\t\t", "URL")
	fmt.Printf("%s\t\t%s\t%s\t%s%s\n", "------", "-----------", "-----------", "\t\t\t", "---")
}
