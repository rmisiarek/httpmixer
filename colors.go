package main

import (
	"fmt"
	"runtime"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[37m"
	white  = "\033[97m"
)

func Red(txt string) string {
	return fmt.Sprintf("%v%v%v", red, txt, reset)
}

func Green(txt string) string {
	return fmt.Sprintf("%v%v%v", green, txt, reset)
}

func Yellow(txt string) string {
	return fmt.Sprintf("%v%v%v", yellow, txt, reset)
}

func Blue(txt string) string {
	return fmt.Sprintf("%v%v%v", blue, txt, reset)
}

func Purple(txt string) string {
	return fmt.Sprintf("%v%v%v", purple, txt, reset)
}

func Cyan(txt string) string {
	return fmt.Sprintf("%v%v%v", cyan, txt, reset)
}

func Gray(txt string) string {
	return fmt.Sprintf("%v%v%v", gray, txt, reset)
}

func White(txt string) string {
	return _color(txt, white, runtime.GOOS)
}

func _color(txt, color, platform string) string {
	if platform == "linux" {
		return fmt.Sprintf("%v%v%v", color, txt, reset)
	} else {
		return txt
	}
}
