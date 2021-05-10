package main

import (
	"fmt"
	"strconv"
	"time"
)

var fmtBase = "%s\t\t%s\t%s"
var fmtFinal = "%-70s\t%s\n"
var fmtSummary = "%s\t\t%s\t%s\n"

func printResult(o *HttpMixerResult) {
	status := strconv.Itoa(o.statusCode)
	description := o.description

	switch string(status[0]) {
	case InformationalCategory:
		p := fmt.Sprintf(fmtBase, Blue(o.method), Blue(status), Blue(description))
		fmt.Printf(fmtFinal, p, White(o.url))
	case SuccessCategory:
		p := fmt.Sprintf(fmtBase, Blue(o.method), Green(status), Green(description))
		fmt.Printf(fmtFinal, p, White(o.url))
	case RedirectionCategory:
		p := fmt.Sprintf(fmtBase, Blue(o.method), Yellow(status), Yellow(description))
		fmt.Printf(fmtFinal, p, White(o.url))
	case ClientErrorCategory:
		p := fmt.Sprintf(fmtBase, Blue(o.method), Purple(status), Purple(description))
		fmt.Printf(fmtFinal, p, White(o.url))
	case ServerErrorCategory:
		p := fmt.Sprintf(fmtBase, Blue(o.method), Red(status), Red(description))
		fmt.Printf(fmtFinal, p, White(o.url))
	default:
		p := fmt.Sprintf(fmtBase, Blue(o.method), Gray(status), Gray(description))
		fmt.Printf(fmtFinal, p, White(o.url))
	}
}

func aggregateSummary(result *HttpMixerResult, showAll bool) {
	if _, exist := summaryData[result.method]; !exist {
		summaryData[result.method] = make(map[string]int)
	}

	label := ""
	if showAll {
		if result.description != "" {
			label = fmt.Sprintf("%d\t%s", result.statusCode, result.description)
		} else {
			label = fmt.Sprintf("%d", result.statusCode)
		}
		summaryData[result.method][label]++
	}

	if !showAll {
		if result.description != "" {
			label = fmt.Sprintf("%d\t%s", result.statusCode, result.description)
			summaryData[result.method][label]++
		}
	}
}

func printSummary(summary Summary, took time.Duration) {
	fmt.Printf(fmtSummary, "\nMETHOD", "FOUND", "RESPONSE STATUS")
	fmt.Printf(fmtSummary, "======", "=====", "===============")
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
			fmt.Printf(fmtSummary, Blue(method), Green(strconv.Itoa(counter)), coloredStatus)
		}
	}

	fmt.Printf("\n>> %s %s\n", Blue("Done within:"), Green(took.String()))
}
