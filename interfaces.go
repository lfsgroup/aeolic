package aeolic

import (
	"io"
	"net/http"
)

// httpClient - client interface
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
