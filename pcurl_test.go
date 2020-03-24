package pcurl

import (
	"fmt"
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
			var gotHeader H
			c.BindHeader(&gotHeader)
			fmt.Printf("==================gotHeader:%s\n", gotHeader)
			if assert.Equal(t, need, gotHeader) {
				c.String(200, "")
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
			curlHeader: []string{"curl", "-X", "POST", "-H", "h1:v1", "-H", "h2:v2"},
			need: H{
				"h1": "v1",
				"h2": "v2",
			},
		},
		testHeader{
			curlHeader: []string{"curl", "-X", "POST", "--header", "h1:v1", "--header", "h2:v2"},
			need: H{
				"h1": "v1",
				"h2": "v2",
			},
		},
	} {

		code := 0
		ts := createGeneralHeader(headerData.need, t)

		req, err := ParseSlice(append(headerData.curlHeader, ts.URL)).Request()
		assert.NoError(t, err)

		gout.New().SetRequest(req).Debug(true).Code(&code).Do()
		assert.Equal(t, code, 200)
	}
}
