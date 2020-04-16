package main

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

func TestGoldilox(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if len(req.Header.Get("junk0")) == 1 {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusRequestHeaderFieldsTooLarge)
		}
	}))

	defer server.Close()

	cmd := exec.Command("./goldilox", "-url", server.URL, "-min", "1", "-max", "2")
	out, err := cmd.CombinedOutput()
	sout := string(out)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !strings.Contains(sout, "Max Header Length: 73") {
		t.Errorf("%v", sout)
	}
}

func TestGoldixAllOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	cmd := exec.Command("./goldilox", "-url", server.URL, "-min", "1", "-max", "2")
	out, err := cmd.CombinedOutput()
	sout := string(out)

	if err == nil {
		t.Errorf("%s", "did not receive exit status as expected")
	}

	if !strings.Contains(sout, "error: all requests succeeded in provided range") {
		t.Errorf("%v", sout)
	}
}

func TestGoldixAllTooLarge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusRequestHeaderFieldsTooLarge)
	}))

	defer server.Close()

	cmd := exec.Command("./goldilox", "-url", server.URL, "-min", "1", "-max", "2")
	out, err := cmd.CombinedOutput()
	sout := string(out)

	if err == nil {
		t.Errorf("%s", "did not receive exit status as expected")
	}

	if !strings.Contains(sout, "error: all requests failed in provided range") {
		t.Errorf("%v", sout)
	}
}
