package aeolic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

// New - returns a new client with the templates loaded from the provided directory.
func New(apiKey string, templateDir string) (Client, error) {
	files, err := os.ReadDir(templateDir)
	if err != nil {
		return Client{}, err
	}
	templates, err := withTemplates(files, ".tmpl.json")
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

// NewWithMap - returns a new client with the provided custom template map.
// Use this method if wish to leverage the embed feature of golang more information can be found: https://pkg.go.dev/embed.
// Note they template name will be the key of the map.

// working example: cmd/embed_slack/main.go
func NewWithMap(apiKey string, templateMap map[string]string) Client {
	return Client{
		Templates:      templateMap,
		DefaultHeaders: setDefaultHeaders(apiKey),
		HTTPClient:     setDefaultClient(),
	}
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

	// hydrate blocks data
	payload := slackChannelPayload{}
	if err := json.Unmarshal(parsedOutput, &payload); err != nil {
		return err
	}

	payload.Channel = channel

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
