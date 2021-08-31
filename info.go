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
	fmt.Printf("%s\n", WhiteBold(title))
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprSource())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprOutput())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprSkipHttps())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprSkipHttp())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprTestTrace())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprPipe())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprRedirect())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprTimeout())
	fmt.Printf("%s %s\n", WhiteBold("=>"), o.reprConcurenncy())
	fmt.Printf("%s %s\n\n", WhiteBold("=>"), o.reprStatusFilter())
}
