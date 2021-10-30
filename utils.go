package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO: refactor intSliceContains & stringSliceContains when Generics will be available

func intSliceContains(s []int, v int) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}

	return false
}

func intArrayToString(a []int) string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}

	return strings.Join(b, ",")
}

func intKeysToSlice(m map[int]string) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func openStdinOrFile(inputs string) (io.ReadCloser, error) {
	r := os.Stdin

	if inputs != "" {
		var err error

		r, err = openFile(inputs)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func openFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createFile(filepath string, sleepSec int) *os.File {
	exist := fileExists(filepath)
	if exist {
		log.Println(Red(fmt.Sprintf(
			"%s exists and will be overwritten. Are you sure? %d seconds to GO\n",
			filepath, sleepSec,
		)))
		time.Sleep(time.Duration(sleepSec) * time.Second)
	}

	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	return file
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
