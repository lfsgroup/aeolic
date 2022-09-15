package aeolic

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// call - creates a new HTTP request and returns an HTTP response
func call(url string, method string, body io.Reader, client httpClient, headers ...map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &http.Response{}, err
	}
	allHeaders := mergeHeaders(headers...)
	for key, value := range allHeaders {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode > 399 {
		return resp, &APIError{
			StatusCode: resp.StatusCode,
			StatusText: http.StatusText(resp.StatusCode),
		}
	}

	if resp.Body == nil {
		return resp, nil
	}

	var slackErr slackErrorResp
	if err := json.NewDecoder(resp.Body).Decode(&slackErr); err != nil {
		return resp, err
	}

	if slackErr.OK {
		return resp, nil
	}

	return resp, &APIError{
		StatusCode: http.StatusBadRequest,
		StatusText: http.StatusText(http.StatusBadRequest),
		Message:    slackErr.Error,
		Context:    errorMessageContextUrl,
	}
}

// mergeHeaders - merge a slice of headers
func mergeHeaders(headersList ...map[string]string) map[string]string {
	mergedHeaders := map[string]string{}
	for _, headers := range headersList {
		for key, value := range headers {
			mergedHeaders[key] = value
		}
	}
	return mergedHeaders
}

// setDefaultHeaders - set default slack headers
func setDefaultHeaders(apiKey string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", apiKey),
		"Content-Type":  "application/json",
	}
}

// setDefaultClient - returns the default http client
func setDefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 15 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: time.Second * 15,
	}
}
