package pcurl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
)

func createGeneralHeader(need H, t *testing.T) *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		cb := func(c *gin.Context) {
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
		}

		router.POST("/", cb)
		router.GET("/", cb)

		return router
	}()

	return httptest.NewServer(http.HandlerFunc(router.ServeHTTP))
}

func Test_Header(t *testing.T) {

	type testHeader struct {
		curlHeader []string
		need       H
	}

	for index, headerData := range []testHeader{
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
		testHeader{
			curlHeader: []string{`curl`, `-H`, `Connection: keep-alive`, `-H`, `Cache-Control: max-age=0`, `-H`, `Upgrade-Insecure-Requests: 1`, `-H`, `Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9`, `-H`, `Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7`, `-H`, `Cookie: username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13aeaa51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; timezone=8; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589004902,1589265192,1589341266,1589717172; Hm_lpvt_12d9f8f1740b76bb88c6691ea1672d8b=1589719525`},
			need: H{
				"Connection":                `keep-alive`,
				"Cache-Control":             `max-age=0`,
				"Upgrade-Insecure-Requests": `1`,
				"Accept":                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9`,
				"Accept-Language":           `en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7`,
				"Cookie":                    `username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13aeaa51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; timezone=8; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589004902,1589265192,1589341266,1589717172; Hm_lpvt_12d9f8f1740b76bb88c6691ea1672d8b=1589719525`,
			},
		},
	} {

		code := 0
		// 创建测试服务端
		ts := createGeneralHeader(headerData.need, t)

		// 解析curl表达式
		req, err := ParseSlice(append(headerData.curlHeader, ts.URL)).Request()
		assert.NoError(t, err, fmt.Sprintf("index = %d", index))

		var getJSON H
		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindJSON(&getJSON).Do()

		assert.NoError(t, err)
		assert.Equal(t, headerData.need, getJSON, fmt.Sprintf("index = %d", index))
		assert.Equal(t, code, 200, fmt.Sprintf("index = %d", index))
	}
}
