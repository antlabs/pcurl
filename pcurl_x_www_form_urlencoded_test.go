package pcurl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
)

type testWWWForm struct {
	Hello string `form:"hello"`
	Pcurl string `form:"pcurl"`
}

func createWWWForm(t *testing.T, need testWWWForm) *httptest.Server {
	r := gin.New()

	r.POST("/", func(c *gin.Context) {
		wf := testWWWForm{}

		var buf bytes.Buffer
		_, err := io.Copy(&buf, c.Request.Body)
		assert.NoError(t, err)

		ioutil.NopCloser(&buf)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		err = c.ShouldBind(&wf)

		assert.NoError(t, err)
		//err := c.ShouldBind(&wf)
		assert.Equal(t, need, wf)
		io.Copy(c.Writer, &buf)
	})

	return httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
}

// --data-urlencode
func Test_DataURLEncode_Option(t *testing.T) {
	type testData struct {
		sendNeed   testWWWForm
		need       string
		curlSlice  []string
		curlString string
		path       string
	}

	sendRequest := func(d testData, index int, req *http.Request) {
		got := ""
		err := gout.New().SetRequest(req).Debug(true).BindBody(&got).Do()
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))
		assert.Equal(t, d.need, got, fmt.Sprintf("test index:%d", index))
	}

	for index, d := range []testData{
		{
			sendNeed:   testWWWForm{Hello: "world", Pcurl: "pcurl"},
			need:       "hello=world&pcurl=pcurl",
			curlSlice:  []string{"curl", "-X", "POST", "--data-urlencode", "hello=world", "--data-urlencode", "pcurl=pcurl"},
			curlString: `curl  -X  POST --data-urlencode hello=world --data-urlencode pcurl=pcurl`,
			path:       "/",
		},
	} {

		// 创建测试服务
		ts := createWWWForm(t, d.sendNeed)

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

		// 发送请求
		sendRequest(d, index, req)
		// =========分隔符================

		// 生成curl字符串
		curlString := d.curlString + " " + url

		fmt.Printf("\nindex:%d#%s\n", index, curlString)

		req, err = ParseString(curlString).Request()
		assert.NoError(t, err, fmt.Sprintf("test index :%d", index))

		// 发送请求
		sendRequest(d, index, req)
	}
}
