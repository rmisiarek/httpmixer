package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client *http.Client

var noRedirect = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func main() {
	source := "debug/google.txt"
	redirect := true

	o := &HttpMixerOpts{
		source:    &source,
		redirect:  &redirect,
		testHttp:  true,
		testHttps: true,
	}

	f := func(o *HttpMixerRes) {
		fmt.Println(o)
	}

	mixer := NewHttpMixer(o)
	mixer.Start(f)
}

func whatStatus(url string) {
	reqURL, reqStatusCode, err := checkStatus(url)
	if err != nil {
		return
	}

	printResult(reqURL, reqStatusCode)
}

func checkStatus(url string) (string, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	return resp.Request.URL.String(), resp.StatusCode, nil
}

func getClient(re *bool) *http.Client {
	client = &http.Client{
		Timeout: 15 * time.Second,
	}

	if *re {
		client.CheckRedirect = noRedirect
	}

	return client
}

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

func prepareUrlsWithSchema(url string) []string {
	result := []string{}

	if !strings.HasPrefix(url, "http") {
		result = append(result, "http://"+url)
	}

	if !strings.HasPrefix(url, "https") {
		result = append(result, "https://"+url)
	}

	return result
}
