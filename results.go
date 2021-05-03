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
	fmt.Println(repeat)
	if ((maxLenOfDescription - descriptionLen) % 8) > 5 {
		return strings.Repeat("\t", int(math.Ceil(repeat)))
	}
	return strings.Repeat("\t", int(math.Floor(repeat)))
}

func printResult(o *HttpMixerResult) {
	s := strconv.Itoa(o.statusCode)

	if description, ok := StatusInformational[o.statusCode]; ok {
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Blue(s), Blue(description), tab(len(description)), o.url)
		return
	}

	if description, ok := StatusSuccess[o.statusCode]; ok {
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Green(s), Green(description), tab(len(description)), o.url)
		return
	}

	if description, ok := StatusRedirection[o.statusCode]; ok {
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Yellow(s), Yellow(description), tab(len(description)), o.url)
		return
	}

	if description, ok := StatusClientError[o.statusCode]; ok {
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Purple(s), Purple(description), tab(len(description)), o.url)
		return
	}

	if description, ok := StatusServerError[o.statusCode]; ok {
		fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Red(s), Red(description), tab(len(description)), o.url)
		return
	}

	fmt.Printf("%s\t\t%s\t\t%s\t%s%s\n", Blue(o.method), Gray(s), Gray("not found"), tab(len("not found")), o.url)
}
