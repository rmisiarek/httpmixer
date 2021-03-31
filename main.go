package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var client *http.Client

var noRedirect = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func main() {
	run()
}

func run() {
	re := true
	client = getClient(&re)

	file := openFile("debug/githubapp.com.txt")
	defer file.Close()

	wg := &sync.WaitGroup{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls := prepareUrlsWithSchema(scanner.Text())
		for _, url := range urls {
			u := url
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				whatStatus(u)
			}(u)
		}
	}
	wg.Wait()
}

func whatStatus(url string) {
	reqURL, reqStatusCode, err := checkStatus(url)
	if err != nil {
		return
	}

	printResult(reqURL, reqStatusCode)
}

// func filterStatus(status int) bool {
// 	return true
// }

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

func openFile(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
