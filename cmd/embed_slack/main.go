package main

import (
	"log"
	"os"

	_ "embed"

	"github.com/lfsgroup/aeolic"
)

//go:embed templates/basic.tmpl.json
var basicTemplate string

func main() {
	// see env variables .env.local
	token := os.Getenv("SLACK_API_TOKEN")
	channel := os.Getenv("TEST_SLACK_CHANNEL")

	customMap := map[string]string{
		"basic": basicTemplate,
	}

	c := aeolic.NewWithMap(token, customMap)

	if err := c.SendMessage(channel, "basic", map[string]string{
		"user_name": "Allan Bond",
	}); err != nil {
		log.Fatal("failed ", err)
	}

}
