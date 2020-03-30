package pcurl

import (
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type H map[string]interface{}

func createGeneral(data string) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {
			if len(data) > 0 {
				c.String(200, data)
			}
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func Test_ParseSlice(t *testing.T) {

	// 创建测试服务
	need := `{"key":"val"}`
	ts := createGeneral(need)
	got := ""
	s := []string{"curl", "-X", "POST", "-d", need, ts.URL}
	req, err := ParseSlice(s).Request()
	assert.NoError(t, err)

	err = gout.New().SetRequest(req).Debug(true).BindBody(&got).Do()
	assert.NoError(t, err)
	assert.Equal(t, need, got)
}

func createGeneralHeader(need H, t *testing.T) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {
			gotHeader := make(H, 2)
			//c.ShouldBindHeader(&gotHeader)

			for k, v := range c.Request.Header {
				if len(v) == 0 {
					continue
				}
				switch k {
				case "Accept-Encoding", "Content-Length", "User-Agent":
					continue
				}

				gotHeader[k] = v[0]
			}

			//c.ShouldBindHeader(&gotHeader2)
			if assert.Equal(t, need, gotHeader) {
				c.JSON(200, gotHeader)
				return
			}

			c.String(500, "")
		})

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func Test_Header(t *testing.T) {

	type testHeader struct {
		curlHeader []string
		need       H
	}

	for _, headerData := range []testHeader{
		testHeader{
			curlHeader: []string{"curl", "-X", "POST", "-H", "H1:v1", "-H", "H2:v2"},
			need: H{
				"H1": "v1",
				"H2": "v2",
			},
		},
		testHeader{
			curlHeader: []string{"curl", "-X", "POST", "--header", "H1:v1", "--header", "H2:v2"},
			need: H{
				"H1": "v1",
				"H2": "v2",
			},
		},
	} {

		code := 0
		// 创建测试服务端
		ts := createGeneralHeader(headerData.need, t)

		// 解析curl表达式
		req, err := ParseSlice(append(headerData.curlHeader, ts.URL)).Request()
		assert.NoError(t, err)

		var getJSON H
		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindJSON(&getJSON).Do()
		assert.NoError(t, err)
		assert.Equal(t, code, 200)
		assert.Equal(t, headerData.need, getJSON)
	}
}

func createGeneralForm(need H, t *testing.T) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.POST("/", func(c *gin.Context) {
			gotForm := make(H, 2)
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

func Test_Form(t *testing.T) {
	type testForm struct {
		curlForm []string
		need     H
	}

	for _, formData := range []testForm{
		testForm{
			curlForm: []string{"curl", "-X", "POST", "-F", "text=good", "-F", "voice=@./testdata/voice.pcm"},
			need: H{
				"text":  "good",
				"voice": "voice",
			},
		},
		testForm{
			curlForm: []string{"curl", "-X", "POST", "--form", "text=good", "--form", "voice=@./testdata/voice.pcm"},
			need: H{
				"text":  "good",
				"voice": "voice",
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
	}
}
