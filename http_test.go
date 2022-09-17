package aeolic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_call_POST_should_not_return_error_and_match_req(t *testing.T) {

	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodPost, nil, m, expectedHeaders)
	assert.NoError(t, err)
}

func Test_call_POST_4xx_should_return_error(t *testing.T) {
	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	var expectedErr *APIError

	_, err := call(url, http.MethodPost, nil, m, expectedHeaders)
	if err == nil {
		t.Error("expected error, got none")
		return
	}
	assert.ErrorAs(t, err, &expectedErr)
}

func Test_call_POST_5xx_should_return_error(t *testing.T) {
	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	var expectedErr *APIError

	_, err := call(url, http.MethodPost, nil, m, expectedHeaders)
	if err == nil {
		t.Error("expected error, got none")
		return
	}
	assert.ErrorAs(t, err, &expectedErr)
}

func Test_call_GET_should_not_return_error_and_match_req(t *testing.T) {
	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodGet, nil, m, expectedHeaders)
	assert.NoError(t, err)

}

func Test_call_PUT_should_not_return_error_and_match_req(t *testing.T) {
	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodPut, nil, m, expectedHeaders)
	assert.NoError(t, err)

}

func Test_call_PATCH_should_not_return_error_and_match_req(t *testing.T) {
	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodPatch, nil, m, expectedHeaders)
	assert.NoError(t, err)

}

func Test_call_DELETE_should_not_return_error_and_match_req(t *testing.T) {
	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodDelete, nil, m, expectedHeaders)
	assert.NoError(t, err)

}

func Test_call_body_should_match(t *testing.T) {
	m := &httpClientMock{}

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}

	body := slackErrorResp{
		OK:    true,
		Error: "",
	}

	data, err := json.Marshal(&body)
	assert.NoError(t, err)

	url := "foo"

	_, err = call(url, http.MethodPost, bytes.NewReader(data), m)
	assert.NoError(t, err)

}

func Test_call_slack_error_should_return_api_error(t *testing.T) {
	m := &httpClientMock{}

	respBody := slackErrorResp{
		OK:    false,
		Error: "invalid_blocks",
	}

	respData, err := json.Marshal(&respBody)
	assert.NoError(t, err)

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respData)),
		}, nil
	}

	body := map[string]string{}

	data, err := json.Marshal(&body)
	assert.NoError(t, err)

	expectedErr := APIError{
		StatusCode: http.StatusBadRequest,
		StatusText: http.StatusText(http.StatusBadRequest),
		Message:    "invalid_blocks",
		Context:    "",
	}
	_, err = call("url", http.MethodPost, bytes.NewReader(data), m)
	assert.Error(t, err)
	assert.Equal(t, expectedErr.Error(), err.Error())

}

func Test_call_slack_ok_should_return_no_error(t *testing.T) {
	m := &httpClientMock{}

	respBody := slackErrorResp{
		OK:    true,
		Error: "",
	}

	respData, err := json.Marshal(&respBody)
	assert.NoError(t, err)

	m.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respData)),
		}, nil
	}

	body := map[string]string{}

	data, err := json.Marshal(&body)
	assert.NoError(t, err)

	_, err = call("url", http.MethodPost, bytes.NewReader(data), m)
	assert.NoError(t, err)

}

func Test_mergeHeaders_should_merge_correctly(t *testing.T) {
	expected := map[string]string{
		"foo": "bar",
		"bin": "baz",
	}

	test := mergeHeaders(map[string]string{"foo": "bar"}, map[string]string{"bin": "baz"})

	assert.Equal(t, expected, test)
}

func Test_mergeHeaders_empty_should_work(t *testing.T) {
	expected := map[string]string{}

	test := mergeHeaders()

	assert.Equal(t, expected, test)
}

func Test_setDefaultHeaders(t *testing.T) {
	want := map[string]string{
		"Authorization": "Bearer 13",
		"Content-Type":  "application/json",
	}

	got := setDefaultHeaders("13")

	assert.Equal(t, want, got)
}
