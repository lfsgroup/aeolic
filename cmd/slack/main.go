package main

import (
	"log"
	"os"

	"github.com/code-gorilla-au/aeolic"
)

func main() {

	token := os.Getenv("SLACK_API_TOKEN")
	channel := os.Getenv("TEST_SLACK_CHANNEL")

	c, err := aeolic.New(token, "./examples")

	if err != nil {
		log.Fatal("error init slack client ", err)
	}

	if err := c.SendMessage(channel, "basic", map[string]string{
		"url_link":  "some-link",
		"user_name": "some-user-name",
	}); err != nil {
		log.Fatal("failed ", err)
	}

}
