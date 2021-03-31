package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
)

var client *http.Client

var withoutRedirect = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func init() {
	client = &http.Client{
		CheckRedirect: withoutRedirect,
		Timeout:       15 * time.Second,
	}
}

func main() {
	file := openFile("debug/githubapp.com.txt")
	defer file.Close()

	wp := workerpool.New(4)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls := prepareUrlsWithSchema(scanner.Text())
		for _, url := range urls {
			u := url
			wp.Submit(func() {
				whatStatus(u)
			})
		}
	}

	wp.StopWait()
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
