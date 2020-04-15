package pcurl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试错误情况
func Test_GetArgsTokenFail(t *testing.T) {
	type TestArgs struct {
		in string
	}

	for _, v := range []TestArgs{
		TestArgs{
			in: `'hello`,
		},
		TestArgs{
			in: `"hello`,
		},
	} {

		_, err := GetArgsToken(v.in)
		assert.Error(t, err)
	}
}

// 测试正确的情况
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
			in:   `curl --location --request DELETE '192.168.5.213:10010/delete/rule?appkey=xx' --header 'Content-Type: text/plain' --data-raw '{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416}'`,
			need: []string{`curl`, `--location`, `--request`, `DELETE`, `192.168.5.213:10010/delete/rule?appkey=xx`, `--header`, `Content-Type: text/plain`, `--data-raw`, `{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416}`},
		},
		TestArgs{
			in:   `'{"s":"{\"s\":\"S\"}"}'`,
			got:  []string{`{"s":"{\"s\":\"S\"}"}`},
			need: []string{`{"s":"{\"s\":\"S\"}"}`},
		},
	} {

		got, err := GetArgsToken(v.in)
		assert.NoError(t, err)
		assert.Equal(t, v.need, got)
	}
}
