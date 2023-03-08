package pcurl

import (
	"fmt"
	"testing"

	"github.com/antlabs/pcurl/body"
	"github.com/stretchr/testify/assert"
)

func TestParseAndObj(t *testing.T) {
	// TODO 如果没有加右'接尾，报错
	all, err := ParseAndObj(`curl https://api.openai.com/v1/completions -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_API_KEY' -d '{ "model": "text-davinci-003", "prompt": "Say this is a test", "max_tokens": 7, "temperature": 0 }'`)
	assert.NoError(t, err)
	assert.Equal(t, all.Encode.Body, body.EncodeJSON)
	assert.Equal(t, all.Method, "POST")
	fmt.Printf("%#v\n", all)
}

func TestParseAndJSON(t *testing.T) {
	// TODO 如果没有加右'接尾，报错
	all, err := ParseAndJSON(`curl https://api.openai.com/v1/completions -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_API_KEY' -d '{ "model": "text-davinci-003", "prompt": "Say this is a test", "max_tokens": 7, "temperature": 0 }'`)
	assert.NoError(t, err)
	fmt.Printf("%s\n", all)
}
