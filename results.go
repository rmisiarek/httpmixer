package main

import (
	"fmt"
	"strconv"
	"time"
)

const fmtBase = "%s %s\t[ %s, %s ]\n"
const fmtSummary = "%s %s of %s [ %s ] \n"

func printResult(o *HttpMixerResult) {
	status := strconv.Itoa(o.statusCode)

	switch string(status[0]) {
	case InformationalCategory:
		fmt.Printf(fmtBase, WhiteBold("==>"), BlueBold(status), Blue(o.method), o.url)
	case SuccessCategory:
		fmt.Printf(fmtBase, WhiteBold("==>"), GreenBold(status), Blue(o.method), o.url)
	case RedirectionCategory:
		fmt.Printf(fmtBase, WhiteBold("==>"), Yellow(status), Blue(o.method), o.url)
	case ClientErrorCategory:
		fmt.Printf(fmtBase, WhiteBold("==>"), MagentaBold(status), Blue(o.method), o.url)
	case ServerErrorCategory:
		fmt.Printf(fmtBase, WhiteBold("==>"), RedBold(status), Blue(o.method), o.url)
	default:
		fmt.Printf(fmtBase, WhiteBold("==>"), WhiteBold(status), Blue(o.method), o.url)
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
		fmt.Printf("%s %s %s\n", WhiteBold("=>"), Blue("summary:"), White("found nothing :("))
	} else {
		fmt.Printf("\n%s %s\n\n", WhiteBold("=>"), Blue("summary:"))
		for method, statuses := range summary {
			for status, counter := range statuses {
				coloredStatus := White(status)
				switch string(status[0]) {
				case InformationalCategory:
					coloredStatus = BlueBold(status)
				case SuccessCategory:
					coloredStatus = GreenBold(status)
				case RedirectionCategory:
					coloredStatus = YellowBold(status)
				case ClientErrorCategory:
					coloredStatus = MagentaBold(status)
				case ServerErrorCategory:
					coloredStatus = RedBold(status)
				}
				fmt.Printf(fmtSummary, WhiteBold("==>"), GreenBold(strconv.Itoa(counter)), coloredStatus, Blue(method))
			}
		}
	}

	fmt.Printf("\n%s %s %s\n", WhiteBold("=>"), Blue("done within:"), Green(took.String()))
}
