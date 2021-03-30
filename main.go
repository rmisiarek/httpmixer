package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
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

	fmt.Println(reqURL, reqStatusCode)
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
