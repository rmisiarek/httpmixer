package main

import (
	"log"
	"strconv"
)

func printResult(url string, status int) {
	s := strconv.Itoa(status)

	if description, ok := StatusInformational[status]; ok {
		log.Printf("[ %s %s ] %s \n", Blue(s), Blue(description), url)
		return
	}

	if description, ok := StatusSuccess[status]; ok {
		log.Printf("[ %s %s ] %s \n", Green(s), Green(description), url)
		return
	}

	if description, ok := StatusRedirection[status]; ok {
		log.Printf("[ %s %s ] %s \n", Yellow(s), Yellow(description), url)
		return
	}

	if description, ok := StatusClientError[status]; ok {
		log.Printf("[ %s %s ] %s \n", Red(s), Red(description), url)
		return
	}

	if description, ok := StatusServerError[status]; ok {
		log.Printf("[ %s %s ] %s \n", Red(s), Red(description), url)
		return
	}

	log.Printf("[ %s ] %s \n", Gray(s), url)
}
