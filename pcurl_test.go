package pcurl

import (
	"github.com/gin-gonic/gin"
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

	ts := createGeneral("")
	s := []string{"curl", "-X", "POST", "-d", `{"key":"val"}`, ts.URL}
	_, err := ParseSlice(s).Request()
	assert.NoError(t, err)
}
