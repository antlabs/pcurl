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
)

type testWWWForm struct {
	Hello string `form:"hello" www-form:"hello"`
	Pcurl string `form:"pcurl" www-form:"pcurl"`
}

func createWWWForm(t *testing.T, need testWWWForm) *httptest.Server {
	r := gin.New()

	r.POST("/", func(c *gin.Context) {
		wf := testWWWForm{}

		var buf bytes.Buffer
		_, err := io.Copy(&buf, c.Request.Body)
		if err != nil {
			t.Fatalf("copy body failed: %v", err)
		}

		ioutil.NopCloser(&buf)

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		err = c.ShouldBind(&wf)
		if err != nil {
			t.Fatalf("ShouldBind failed: %v", err)
		}
		//err := c.ShouldBind(&wf)
		if wf != need {
			t.Fatalf("unexpected form, got=%+v want=%+v", wf, need)
		}
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
		if err != nil {
			t.Fatalf("request failed (index=%d): %v", index, err)
		}
		if got != d.need {
			t.Fatalf("unexpected body (index=%d): got=%q want=%q", index, got, d.need)
		}
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
		if err != nil {
			t.Fatalf("ParseSlice.Request failed (index=%d): %v", index, err)
		}

		// 发送请求
		sendRequest(d, index, req)
		// =========分隔符================

		// 生成curl字符串
		curlString := d.curlString + " " + url

		fmt.Printf("\nindex:%d#%s\n", index, curlString)

		req, err = ParseString(curlString).Request()
		if err != nil {
			t.Fatalf("ParseString.Request failed (index=%d): %v", index, err)
		}

		// 发送请求
		sendRequest(d, index, req)
	}
}
