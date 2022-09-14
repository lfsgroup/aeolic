package aeolic

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	postChatMessageEndpoint = "https://slack.com/api/chat.postMessage"
)

type Client struct {
	DefaultHeaders map[string]string
	Templates      map[string]string
	HTTPClient     httpClient
}

func New(apiKey string, templateDir string) (Client, error) {
	templates, err := withTemplates(templateDir, ".tmpl")
	if err != nil {
		return Client{}, err
	}
	c := Client{
		Templates:      templates,
		DefaultHeaders: setDefaultHeaders(apiKey),
		HTTPClient:     setDefaultClient(),
	}
	return c, nil
}

func (c *Client) SendMessage(channel string, templateName string, body any) error {

	parsedOutput, err := parse(templateName, c.Templates, body)
	if err != nil {
		return fmt.Errorf("could not parse template [%s] error [%w]", templateName, err)
	}
	if _, err := call(postChatMessageEndpoint, http.MethodPost, bytes.NewReader(parsedOutput), c.HTTPClient, c.DefaultHeaders); err != nil {
		return err
	}
	return nil
}
