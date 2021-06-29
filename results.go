package main

import (
	"fmt"
	"strconv"
	"time"
)

var fmtBase = "==> %s\t[ %s, %s ]\n"

func printResult(o *HttpMixerResult) {
	status := strconv.Itoa(o.statusCode)

	switch string(status[0]) {
	case InformationalCategory:
		fmt.Printf(fmtBase, Blue(status), Blue(o.method), o.url)
	case SuccessCategory:
		fmt.Printf(fmtBase, Green(status), Blue(o.method), o.url)
	case RedirectionCategory:
		fmt.Printf(fmtBase, Yellow(status), Blue(o.method), o.url)
	case ClientErrorCategory:
		fmt.Printf(fmtBase, Purple(status), Blue(o.method), o.url)
	case ServerErrorCategory:
		fmt.Printf(fmtBase, Red(status), Blue(o.method), o.url)
	default:
		fmt.Printf(fmtBase, Gray(status), Blue(o.method), o.url)
	}
}

func aggregateSummary(result *HttpMixerResult, showAll bool) {
	if _, exist := summaryData[result.method]; !exist {
		summaryData[result.method] = make(map[string]int)
	}

	label := ""
	if showAll {
		if result.description != "" {
			label = fmt.Sprintf("%d - %s", result.statusCode, result.description)
		} else {
			label = fmt.Sprintf("%d", result.statusCode)
		}
		summaryData[result.method][label]++
	}

	if !showAll {
		if result.description != "" {
			label = fmt.Sprintf("%d - %s", result.statusCode, result.description)
			summaryData[result.method][label]++
		}
	}
}

func printSummary(summary Summary, took time.Duration) {
	if len(summary) == 0 {
		fmt.Printf(">> %s %s\n", Blue("summary:"), White("found nothing :("))
	} else {
		fmt.Printf("\n>> %s\n\n", Blue("summary:"))
		for method, statuses := range summary {
			for status, counter := range statuses {
				coloredStatus := White(status)
				switch string(status[0]) {
				case InformationalCategory:
					coloredStatus = Blue(status)
				case SuccessCategory:
					coloredStatus = Green(status)
				case RedirectionCategory:
					coloredStatus = Yellow(status)
				case ClientErrorCategory:
					coloredStatus = Purple(status)
				case ServerErrorCategory:
					coloredStatus = Red(status)
				}
				fmt.Printf("==> found: %s of %s [ %s ] \n", Green(strconv.Itoa(counter)), coloredStatus, Blue(method))
			}
		}
	}

	fmt.Printf("\n>> %s %s\n", Blue("done within:"), Green(took.String()))
}
