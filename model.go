package pcurl

import "errors"

const (
	EncodeWWWForm = "www-form"
	EncodeJSON    = "json"
	EncodeYAML    = "yaml"
)

var ErrUnknownEncode = errors.New("Unknown encoder")

type Encode struct {
	// x-www-form-urlencoded
	Body string `json:"body,omitempty" yaml:"body"`
}

type Req struct {
	Method string   `json:"method,omitempty" yaml:"method"`
	URL    string   `json:"url,omitempty" yaml:"url"`
	Encode Encode   `json:"encode,omitempty" yaml:"encode"`
	Body   any      `json:"body,omitempty" yaml:"body"`
	Header []string `json:"header,omitempty" yaml:"header"`
}
