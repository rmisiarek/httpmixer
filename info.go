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
	fmt.Printf("%s\n", Blue(title))
	fmt.Printf(">> %s\n", o.reprSource())
	fmt.Printf(">> %s\n", o.reprOutput())
	fmt.Printf(">> %s\n", o.reprSkipHttps())
	fmt.Printf(">> %s\n", o.reprSkipHttp())
	fmt.Printf(">> %s\n", o.reprTestTrace())
	fmt.Printf(">> %s\n", o.reprRedirect())
	fmt.Printf(">> %s\n", o.reprTimeout())
	fmt.Printf(">> %s\n", o.reprConcurenncy())
	fmt.Printf(">> %s\n\n", o.reprStatusFilter())
	fmt.Printf("%s\t\t%s\t\t\t\t\t%s", "METHOD", "RESPONSE STATUS", "URL\n")
	fmt.Printf("%s\t\t%s\t\t\t\t\t%s", "======", "===============", "===\n")
}
