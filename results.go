package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// TODO: find better way for this
func tab(descriptionLen int) string {
	maxLenOfDescription := 34
	repeat := float64(maxLenOfDescription-descriptionLen) / float64(8)
	if ((maxLenOfDescription - descriptionLen) % 8) >= 5 {
		return strings.Repeat("\t", int(math.Ceil(repeat)))
	}
	return strings.Repeat("\t", int(math.Floor(repeat)))
}

func printResult(o *HttpMixerResult) {
	status := strconv.Itoa(o.statusCode)
	description := o.resolvedCategory.description

	switch o.resolvedCategory.category {
	case InformationalCategory:
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Blue(status), Blue(description), tab(len(description)), o.url)
	case SuccessCategory:
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Green(status), Green(description), tab(len(description)), o.url)
	case RedirectionCategory:
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Yellow(status), Yellow(description), tab(len(description)), o.url)
	case ClientErrorCategory:
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Purple(status), Purple(description), tab(len(description)), o.url)
	case ServerErrorCategory:
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Red(status), Red(description), tab(len(description)), o.url)
	default:
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Gray(status), Gray("not found"), tab(len("not found")), o.url)
	}
}
