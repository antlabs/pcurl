package pcurl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetArgsToken(t *testing.T) {
	type TestArgs struct {
		in   string
		got  []string
		need []string
	}

	for _, v := range []TestArgs{
		TestArgs{
			in:   `curl -XGET "http://192.168.6.100:9200/eval-log/_search" -H 'Content-Type: application/json' -d'{  "query": {    "match": {      "level": "error"    }  }}'`,
			need: []string{`curl`, `-XGET`, "http://192.168.6.100:9200/eval-log/_search", "-H", `Content-Type: application/json`, `-d{  "query": {    "match": {      "level": "error"    }  }}`},
		},
		TestArgs{
			in:   `curl --location --request DELETE '192.168.5.213:10010/delete/rule?appkey=xx' --header 'Content-Type: text/plain' --data-raw '{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416`,
			need: []string{`curl`, `--location`, `--request`, `DELETE`, `192.168.5.213:10010/delete/rule?appkey=xx`, `--header`, `Content-Type: text/plain`, `--data-raw`, `{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416`},
		},
	} {

		got := GetArgsToken(v.in)
		assert.Equal(t, v.need, got)
	}
}
