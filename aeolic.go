package aeolic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	postChatMessageEndpoint = "https://slack.com/api/chat.postMessage"
	errorMessageContextUrl  = "https://api.slack.com/methods/chat.postMessage#errors"
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

type slackChannelPayload struct {
	Channel string          `json:"channel,omitempty"`
	Blocks  json.RawMessage `json:"blocks,omitempty"`
}

// SendMessage - post a slack message to a channel
func (c *Client) SendMessage(channel string, templateName string, body any) error {

	parsedOutput, err := parse(templateName, c.Templates, body)
	if err != nil {
		return fmt.Errorf("could not parse template [%s] error [%w]", templateName, err)
	}

	payload := slackChannelPayload{
		Channel: channel,
		Blocks:  parsedOutput,
	}

	data, err := json.Marshal(&payload)

	if err != nil {
		return err
	}

	_, err = call(postChatMessageEndpoint, http.MethodPost, bytes.NewReader(data), c.HTTPClient, c.DefaultHeaders)
	if err != nil {
		return err
	}
	return nil
}
