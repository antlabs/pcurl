package pcurl

import (
	"fmt"
	"testing"

	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
)

func Test_Data(t *testing.T) {

	type testData struct {
		need       string
		curlSlice  []string
		curlString string
		path       string
	}

	for index, d := range []testData{
		{
			need:       `{"key":"val"}`,
			curlSlice:  []string{"curl", "-X", "POST", "-d", `{"key":"val"}`},
			curlString: `curl  -X  POST -d '{"key":"val"}'`,
			path:       "/",
		},
		{ //测试-k选项
			need:       `{"key":"val"}`,
			curlSlice:  []string{"curl", "-k", "-X", "POST", "-d", `{"key":"val"}`},
			curlString: `curl -k -X  POST -d '{"key":"val"}'`,
			path:       "/",
		},
		{ //测试--insecure选项
			need:       `{"key":"val"}`,
			curlSlice:  []string{"curl", "--insecure", "-X", "POST", "-d", `{"key":"val"}`},
			curlString: `curl --insecure -X  POST -d '{"key":"val"}'`,
			path:       "/",
		},
		{
			need:       `{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416"}`,
			curlSlice:  []string{"curl", "--location", "--request", "DELETE", "--header", "Content-Type: text/plain", "--data-raw", `{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416"}`},
			curlString: `curl --location --request DELETE --header 'Content-Type: text/plain' --data-raw '{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416"}'`,
			path:       "/appkey/admin/v1/delete/connect/rule?appkey=xx",
		},
		{
			need:       `{"type":"region","list":[{"region":"sh","business":"asr","protocol":"http","connect":56},{"region":"bj","business":"asr","protocol":"websocket","connect":52},{"region":"bj","business":"asr","protocol":"private","connect":51}]}`,
			curlSlice:  []string{"curl", "--location", "--request", "POST", "--header", "Content-Type: text/plain", "--data-raw", `{"type":"region","list":[{"region":"sh","business":"asr","protocol":"http","connect":56},{"region":"bj","business":"asr","protocol":"websocket","connect":52},{"region":"bj","business":"asr","protocol":"private","connect":51}]}`},
			curlString: `curl --location --request POST --header 'Content-Type: text/plain' --data-raw '{"type":"region","list":[{"region":"sh","business":"asr","protocol":"http","connect":56},{"region":"bj","business":"asr","protocol":"websocket","connect":52},{"region":"bj","business":"asr","protocol":"private","connect":51}]}'`,
			path:       "/appkey/admin/v1/add/connect/rule?appkey=xx",
		},
	} {

		// 创建测试服务
		ts := createGeneral2()
		got := ""

		// 生成curl slice
		url := ts.URL
		if len(d.path) > 0 {
			url = url + d.path
		}

		// curlSlice追加url
		curlSlice := append(d.curlSlice, url)

		fmt.Printf("\nindex:%d#%s\n", index, curlSlice)

		req, err := ParseSlice(curlSlice).Request()
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))

		err = gout.New().SetRequest(req).Debug(true).BindBody(&got).Do()
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))
		assert.Equal(t, d.need, got)

		// 生成curl字符串
		curlString := d.curlString + " " + url
		fmt.Printf("\nindex:%d#%s\n", index, curlString)
		req, err = ParseString(curlString).Request()
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))

		err = gout.New().SetRequest(req).Debug(true).BindBody(&got).Do()
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))
		assert.Equal(t, d.need, got, fmt.Sprintf("test index:%d", index))
	}

}
