package aeolic

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Client struct {
	apiKey    string
	templates map[string]string
}

func New(apiKey string, templateDir string) (Client, error) {
	templates, err := withTemplates(templateDir, ".tmpl")
	if err != nil {
		return Client{}, err
	}
	return Client{
		apiKey:    apiKey,
		templates: templates,
	}, nil
}

func (c *Client) SendMessage(channel string, templateName string, body any) error {
	return fmt.Errorf("not implemented")
}

func parse(templateName string, templateMap map[string]string, data any) ([]byte, error) {

	slackBlock, ok := templateMap[templateName]
	if !ok {
		return []byte{}, fmt.Errorf("template not found")
	}

	tmpl, err := template.New(templateName).Option("missingkey=error").Parse(slackBlock)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)
	if err := tmpl.Execute(wr, data); err != nil {
		return nil, fmt.Errorf("%w \n %s", err, slackBlock)
	}
	if err := wr.Flush(); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// withTemplates - load templates by file path
func withTemplates(dirPath string, fileSuffix string) (map[string]string, error) {
	rootTemplates := map[string]string{}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return rootTemplates, err
	}

	for _, file := range files {
		fileLocation := filepath.Join(dirPath, file.Name())
		if strings.HasSuffix(file.Name(), fileSuffix) {
			data, err := os.ReadFile(filepath.Clean(fileLocation))
			if err != nil {
				return rootTemplates, err
			}
			rootTemplates[fileLocation] = string(data)
		}
	}
	return rootTemplates, nil
}
