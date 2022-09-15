package aeolic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_call_POST_should_not_return_error_and_match_req(t *testing.T) {
	m := MockHTTPClient{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
		},
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodPost, nil, &m, expectedHeaders)
	assert.NoError(t, err)
	for key, value := range expectedHeaders {
		assert.Equal(t, m.Req.Header.Get(key), value)
	}
	assert.Equal(t, m.Req.URL.String(), url)
	assert.Equal(t, m.Req.Method, http.MethodPost)

}

func Test_call_POST_4xx_should_return_error(t *testing.T) {
	m := MockHTTPClient{
		Resp: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	var expectedErr *APIError

	_, err := call(url, http.MethodPost, nil, &m, expectedHeaders)
	if err == nil {
		t.Error("expected error, got none")
		return
	}
	assert.ErrorAs(t, err, &expectedErr)
}

func Test_call_POST_5xx_should_return_error(t *testing.T) {
	m := MockHTTPClient{
		Resp: &http.Response{
			StatusCode: http.StatusInternalServerError,
		},
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	var expectedErr *APIError

	_, err := call(url, http.MethodPost, nil, &m, expectedHeaders)
	if err == nil {
		t.Error("expected error, got none")
		return
	}
	assert.ErrorAs(t, err, &expectedErr)
}

func Test_call_GET_should_not_return_error_and_match_req(t *testing.T) {
	m := MockHTTPClient{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
		},
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodGet, nil, &m, expectedHeaders)
	assert.NoError(t, err)
	for key, value := range expectedHeaders {
		assert.Equal(t, m.Req.Header.Get(key), value)
	}
	assert.Equal(t, m.Req.URL.String(), url)
	assert.Equal(t, m.Req.Method, http.MethodGet)

}
func Test_call_PUT_should_not_return_error_and_match_req(t *testing.T) {
	m := MockHTTPClient{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
		},
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodPut, nil, &m, expectedHeaders)
	assert.NoError(t, err)
	for key, value := range expectedHeaders {
		assert.Equal(t, m.Req.Header.Get(key), value)
	}
	assert.Equal(t, m.Req.URL.String(), url)
	assert.Equal(t, m.Req.Method, http.MethodPut)

}

func Test_call_PATCH_should_not_return_error_and_match_req(t *testing.T) {
	m := MockHTTPClient{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
		},
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodPatch, nil, &m, expectedHeaders)
	assert.NoError(t, err)
	for key, value := range expectedHeaders {
		assert.Equal(t, m.Req.Header.Get(key), value)
	}
	assert.Equal(t, m.Req.URL.String(), url)
	assert.Equal(t, m.Req.Method, http.MethodPatch)

}

func Test_call_DELETE_should_not_return_error_and_match_req(t *testing.T) {
	m := MockHTTPClient{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
		},
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodDelete, nil, &m, expectedHeaders)
	assert.NoError(t, err)
	for key, value := range expectedHeaders {
		assert.Equal(t, m.Req.Header.Get(key), value)
	}
	assert.Equal(t, m.Req.URL.String(), url)
	assert.Equal(t, m.Req.Method, http.MethodDelete)

}

func Test_call_body_should_match(t *testing.T) {
	m := MockHTTPClient{}

	body := slackErrorResp{
		OK:    true,
		Error: "",
	}

	data, err := json.Marshal(&body)
	assert.NoError(t, err)

	url := "foo"

	_, err = call(url, http.MethodPost, bytes.NewReader(data), &m)
	assert.NoError(t, err)

	test := slackErrorResp{}
	err = json.NewDecoder(m.Req.Body).Decode(&test)
	assert.NoError(t, err)
	assert.Equal(t, test, body)

}

func Test_call_slack_error_should_return_api_error(t *testing.T) {
	m := MockHTTPClient{}

	respBody := slackErrorResp{
		OK:    false,
		Error: "invalid_blocks",
	}

	respData, err := json.Marshal(&respBody)
	assert.NoError(t, err)

	m.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respData)),
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
	_, err = call("url", http.MethodPost, bytes.NewReader(data), &m)
	assert.Error(t, err)
	assert.Equal(t, expectedErr.Error(), err.Error())

}

func Test_call_slack_ok_should_return_no_error(t *testing.T) {
	m := MockHTTPClient{}

	respBody := slackErrorResp{
		OK:    true,
		Error: "",
	}

	respData, err := json.Marshal(&respBody)
	assert.NoError(t, err)

	m.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respData)),
	}

	body := map[string]string{}

	data, err := json.Marshal(&body)
	assert.NoError(t, err)

	_, err = call("url", http.MethodPost, bytes.NewReader(data), &m)
	assert.NoError(t, err)

}

func Test_call_should_should_return_error(t *testing.T) {
	m := MockHTTPClient{
		ErrDo: true,
		Err:   errors.New("expected error"),
	}

	expectedHeaders := map[string]string{
		"Auth": "/app/json",
	}

	url := "foo"

	_, err := call(url, http.MethodPost, nil, &m, expectedHeaders)
	assert.ErrorIs(t, err, m.Err)
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
