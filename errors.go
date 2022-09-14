package aeolic

import "fmt"

type APIError struct {
	StatusCode int
	StatusText string
	Message    string
	Context    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s [%d]: %s", e.StatusText, e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error {
	return fmt.Errorf("%s: [%d]: %s", e.StatusText, e.StatusCode, e.Message)
}

// Slack error response to map
type slackErrorResp struct {
	OK    bool   `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
}
