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
	} {

		got := GetArgsToken(v.in)
		assert.Equal(t, v.need, got)
	}
}
