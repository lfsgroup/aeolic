package aeolic

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse_simple(t *testing.T) {
	got, err := parse("basic", map[string]string{
		"basic": `{ "hello": "{{ .hello }}" }`,
	}, map[string]string{
		"hello": "world",
	})

	assert.NoError(t, err)

	assert.Equal(t, `{ "hello": "world" }`, string(got))
}

//go:embed examples/basic.tmpl.json
var embedTest string

func Test_parse_embed(t *testing.T) {

	got, err := parse("basic", map[string]string{
		"basic": embedTest,
	}, map[string]string{
		"user_name": "some-user-name",
	})

	assert.NoError(t, err)

	assert.Contains(t, string(got), "some-user-name")
}

func Test_parse_missing_key(t *testing.T) {
	_, err := parse("basic", map[string]string{
		"basic": embedTest,
	}, map[string]string{
		"url_link": "some-link",
	})

	assert.Contains(t, err.Error(), "no entry for key \"user_name\" ")

}
