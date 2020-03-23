package pcurl

import (
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
