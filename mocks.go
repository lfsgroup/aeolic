package aeolic

import (
	"io"
	"net/http"
)

type MockHTTPClient struct {
	Retries int
	Resp    *http.Response
	Req     *http.Request
	Err     error
	ErrDo   bool
	ErrGet  bool
	ErrPost bool
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.Retries++
	m.Req = req
	if m.ErrDo {
		return nil, m.Err
	}
	if m.Resp == nil {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}
	return m.Resp, nil
}
func (m *MockHTTPClient) Get(url string) (resp *http.Response, err error) {
	m.Retries++
	if m.ErrGet {
		return nil, m.Err
	}
	if m.Resp == nil {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}
	return m.Resp, nil
}
func (m *MockHTTPClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	m.Retries++
	if m.ErrPost {
		return nil, m.Err
	}
	if m.Resp == nil {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}
	return m.Resp, nil
}
