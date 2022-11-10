package aeolic

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func parse(templateName string, templateMap map[string]string, data any) ([]byte, error) {

	slackTemplates, ok := templateMap[templateName]
	if !ok {
		return []byte{}, fmt.Errorf("template %s does not exist", templateName)
	}

	tmpl, err := template.New(templateName).Option("missingkey=error").Parse(slackTemplates)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)
	if err := tmpl.Execute(wr, data); err != nil {
		return nil, fmt.Errorf("%w \n %s", err, slackTemplates)
	}
	if err := wr.Flush(); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// withTemplates - load templates by directory entry
func withTemplates(files []fs.DirEntry, dirPath, fileSuffix string) (map[string]string, error) {
	rootTemplates := map[string]string{}

	for _, file := range files {
		fileLocation := filepath.Join(dirPath, file.Name())
		if strings.HasSuffix(file.Name(), fileSuffix) {
			data, err := os.ReadFile(fileLocation)
			if err != nil {
				return rootTemplates, err
			}
			stripedFileName := strings.TrimRight(file.Name(), fileSuffix)
			rootTemplates[stripedFileName] = string(data)
		}
	}
	return rootTemplates, nil
}
