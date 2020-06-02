package pcurl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type H map[string]interface{}

func createGeneral(data string) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.Any("/*test", func(c *gin.Context) {
			if len(data) > 0 {
				c.String(200, data)
			}
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func Test_Curl(t *testing.T) {

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
			need:       `{"type":"region","list":[{"region":"sh","business":"asr","protocol":"http","connect":56},{"region":"bj","business":"asr","protocol":"websocket","connect":52},{"region":"bj","business":"asr","protocol":"private","connect":51}]}'`,
			curlSlice:  []string{"curl", "--location", "--request", "POST", "--header", "Content-Type: text/plain", "--data-raw", `{"type":"region","list":[{"region":"sh","business":"asr","protocol":"http","connect":56},{"region":"bj","business":"asr","protocol":"websocket","connect":52},{"region":"bj","business":"asr","protocol":"private","connect":51}]}'`},
			curlString: `curl --location --request POST --header 'Content-Type: text/plain' --data-raw '{"type":"region","list":[{"region":"sh","business":"asr","protocol":"http","connect":56},{"region":"bj","business":"asr","protocol":"websocket","connect":52},{"region":"bj","business":"asr","protocol":"private","connect":51}]}'`,
			path:       "/appkey/admin/v1/add/connect/rule?appkey=xx",
		},
	} {

		// 创建测试服务
		ts := createGeneral(d.need)
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
		assert.Equal(t, d.need, got)
	}

}

// 测试formdata
func Test_Form(t *testing.T) {
	type testForm struct {
		curlForm       []string
		curlFormString string
		need           H
	}

	for _, formData := range []testForm{
		testForm{
			curlForm:       []string{"curl", "-X", "POST", "-F", "text=good", "-F", "voice=@./testdata/voice.pcm"},
			curlFormString: `curl -X POST -F text=good -F voice=@./testdata/voice.pcm`,
			need: H{
				"text":  "good",
				"voice": "voice\n",
			},
		},
		testForm{
			curlForm:       []string{"curl", "--request", "POST", "--form", "text=good", "--form", "voice=@./testdata/voice.pcm"},
			curlFormString: `curl --request POST --form text=good --form voice=@./testdata/voice.pcm`,
			need: H{
				"text":  "good",
				"voice": "voice\n",
			},
		},
	} {

		code := 0
		// 创建测试服务端
		ts := createGeneralForm(formData.need, t)

		// 解析curl表达式
		req, err := ParseSlice(append(formData.curlForm, ts.URL)).Request()
		assert.NoError(t, err)

		var getJSON H
		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindJSON(&getJSON).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
		assert.Equal(t, formData.need, getJSON)

		// 测试string方式
		req, err = ParseAndRequest(formData.curlFormString + " " + ts.URL)
		assert.NoError(t, err)

		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindJSON(&getJSON).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
		assert.Equal(t, formData.need, getJSON)
	}
}

func createGeneralForm(need H, t *testing.T) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {
			gotForm := make(H, 2)
			err := c.Request.ParseMultipartForm(32 * 1024 * 1024)

			assert.NoError(t, err)

			for k, f := range c.Request.Form {
				if len(f) == 0 {
					continue
				}

				gotForm[k] = f[0]
			}
			for k, f := range c.Request.MultipartForm.File {
				if len(f) == 0 {
					continue
				}
				fd, err := f[0].Open()
				assert.NoError(t, err)

				var s strings.Builder

				io.Copy(&s, fd)
				gotForm[k] = s.String()
				fd.Close()
			}

			c.ShouldBind(&gotForm)

			//c.ShouldBindHeader(&gotHeader2)
			if assert.Equal(t, need, gotForm) {
				c.JSON(200, gotForm)
				return
			}

			c.String(500, "")
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func Test_Method(t *testing.T) {
	methodServer := func() *httptest.Server {
		router := func() *gin.Engine {
			router := gin.New()

			router.DELETE("/", func(c *gin.Context) {
				c.String(200, "DELETE")
			})

			router.GET("/", func(c *gin.Context) {
				c.String(200, "GET")
			})
			return router
		}()

		return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
	}

	need := []string{"DELETE", "DELETE", "GET"}

	for index, curlStr := range []string{
		`curl -X DELETE -G `,
		`curl -G -XDELETE `,
		`curl -G `,
	} {
		ts := methodServer()
		req, err := ParseAndRequest(curlStr + ts.URL)

		assert.NoError(t, err)

		got := ""
		err = gout.New().SetRequest(req).BindBody(&got).Do()

		assert.Equal(t, got, need[index])
		assert.NoError(t, err)
	}
}

func Test_URL(t *testing.T) {
	type testURL struct {
		curl []string
		need string
	}

	for _, urlData := range []testURL{
		testURL{
			curl: []string{"curl", "-X", "POST"},
			need: "1",
		},
		testURL{
			curl: []string{"curl", "-X", "POST"},
			need: "2",
		},
	} {

		code := 0
		// 创建测试服务端
		ts := createGeneral("1")
		ts2 := createGeneral("2")

		// 解析curl表达式
		var curl []string
		if urlData.need == "1" {
			curl = append(urlData.curl, "--url", ts2.URL, ts.URL)
		} else {
			curl = append(urlData.curl, ts.URL, "--url", ts2.URL)

		}

		req, err := ParseSlice(curl).Request()
		assert.NoError(t, err)

		s := ""
		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindBody(&s).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
		assert.Equal(t, urlData.need, s)
	}
}
