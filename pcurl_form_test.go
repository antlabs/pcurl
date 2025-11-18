package pcurl

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
)

func createGeneralForm(need H, t *testing.T) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {
			gotForm := make(H, 2)
			err := c.Request.ParseMultipartForm(32 * 1024 * 1024)
			if err != nil {
				t.Errorf("ParseMultipartForm failed: %v", err)
				return
			}

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
				if err != nil {
					t.Errorf("open multipart file failed: %v", err)
					return
				}

				var s strings.Builder

				io.Copy(&s, fd)
				gotForm[k] = s.String()
				fd.Close()
			}

			c.ShouldBind(&gotForm)

			//c.ShouldBindHeader(&gotHeader2)
			if !mapsEqual(need, gotForm) {
				c.String(500, "")
				return
			}
			c.JSON(200, gotForm)
			return
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
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
		if err != nil {
			t.Fatalf("ParseSlice.Request failed: %v", err)
		}

		var getJSON H
		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindJSON(&getJSON).Do()
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if code != 200 {
			t.Fatalf("unexpected status code, got=%d want=%d", code, 200)
		}
		if !mapsEqual(formData.need, getJSON) {
			t.Fatalf("unexpected JSON, got=%v want=%v", getJSON, formData.need)
		}

		// 测试string方式
		req, err = ParseAndRequest(formData.curlFormString + " " + ts.URL)
		if err != nil {
			t.Fatalf("ParseAndRequest failed: %v", err)
		}

		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindJSON(&getJSON).Do()
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if code != 200 {
			t.Fatalf("unexpected status code, got=%d want=%d", code, 200)
		}
		if !mapsEqual(formData.need, getJSON) {
			t.Fatalf("unexpected JSON, got=%v want=%v", getJSON, formData.need)
		}
	}
}
