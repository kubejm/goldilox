package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func reqLength(req *http.Request) int {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	return len(dump)
}

func invoke(url string, hdrLength int) int {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("junk", strings.Repeat("a", hdrLength))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	return resp.StatusCode
}

func search(url string, low int, high int) (int, error) {
	for (high - low) != 1 {
		median := (low + high) / 2

		statusCode := invoke(url, median)
		okStatus := statusCode >= 200 && statusCode < 300

		if okStatus {
			low = median
		} else {
			high = median
		}
	}

	lowStatusCode := invoke(url, low)
	lowStatusOk := lowStatusCode >= 200 && lowStatusCode < 300

	if !lowStatusOk {
		return 0, errors.New("error: all requests failed in provided range")
	}

	highStatusCode := invoke(url, high)
	highStatusOk := highStatusCode >= 200 && highStatusCode < 300

	if highStatusOk {
		return 0, errors.New("error: all requests succeeded in provided range")
	}

	genHdrLength := 50
	urlLength := len(url)
	maxHdrLength := genHdrLength + urlLength + low

	return maxHdrLength, nil
}

func main() {
	url := flag.String("url", "", "url to send GET requests to (required)")
	min := flag.Int("min", 1, "minimum header size (bytes) for tesing range (defaults to 1)")
	max := flag.Int("max", 10000, "maximum header size (bytes) for tesing range (defaults to 10000)")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "url is required, use the -h flag for help")
		os.Exit(1)
	}

	if *min >= *max {
		fmt.Fprintln(os.Stderr, "min cannot be greater than or equal to max")
		os.Exit(1)
	}

	maxHdrLength, err := search(*url, *min, *max)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Max Header Length: %d bytes\n", maxHdrLength)
}
