package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

func chunkString(s string, size int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += size {
		end := i + size

		if end > len(runes) {
			end = len(runes)
		}

		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks
}

func invoke(url string, hdrLength int, chunkSize int) (*http.Request, *http.Response) {
	req, err := http.NewRequest("GET", url, nil)

	chunks := chunkString(strings.Repeat("a", hdrLength), chunkSize)
	for index, chunk := range chunks {
		req.Header.Set("junk"+strconv.Itoa(index), chunk)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	return req, resp
}

func hdrLength(req *http.Request) int {
	data, err := httputil.DumpRequest(req, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	return len(data)
}

func search(url string, low int, high int, chunkSize int) (int, error) {
	for (high - low) != 1 {
		median := (low + high) / 2

		_, resp := invoke(url, median, chunkSize)
		okStatus := resp.StatusCode >= 200 && resp.StatusCode < 300

		if okStatus {
			low = median
		} else {
			high = median
		}
	}

	lowReq, lowResp := invoke(url, low, chunkSize)
	lowStatusOk := lowResp.StatusCode >= 200 && lowResp.StatusCode < 300

	if !lowStatusOk {
		return 0, errors.New("error: all requests failed in provided range")
	}

	_, highResp := invoke(url, high, chunkSize)
	highStatusOk := highResp.StatusCode >= 200 && highResp.StatusCode < 300

	if highStatusOk {
		return 0, errors.New("error: all requests succeeded in provided range")
	}

	maxHdrLength := hdrLength(lowReq)

	return maxHdrLength, nil
}

func main() {
	url := flag.String("url", "", "url to send GET requests to (required)")
	min := flag.Int("min", 1, "minimum header size (bytes) for tesing range")
	max := flag.Int("max", 10000, "maximum header size (bytes) for tesing range")
	chunkSize := flag.Int("chunkSize", 3000, "size (bytes) to partition header into separate key/value paris")
	flag.Parse()

	if *url == "" {
		fmt.Fprintln(os.Stderr, "url is required, use the -h flag for help")
		os.Exit(1)
	}

	if *min >= *max {
		fmt.Fprintln(os.Stderr, "min cannot be greater than or equal to max")
		os.Exit(1)
	}

	maxHdrLength, err := search(*url, *min, *max, *chunkSize)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Max Header Length: %d bytes\n", maxHdrLength)
}
