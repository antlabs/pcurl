package pcurl

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
)

func createCompressed() *httptest.Server {

	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	cb := func(c *gin.Context) {
		//c.Writer.Header().Add("Content-Encoding", "gzip")
		io.Copy(c.Writer, c.Request.Body)
	}

	r.POST("/", cb)

	return httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
}

func Test_Compressed(t *testing.T) {
	router := createCompressed()
	type rspHeader struct {
		ContentEncoding string `header:"Content-Encoding"`
	}

	for _, v := range []string{
		`curl --compressed -d "test compressed test compressed"`,
	} {
		curlString := v + " " + router.URL + "/" + " " + "-d " + strings.Repeat("test", 100)

		req, err := ParseAndRequest(curlString)

		assert.NoError(t, err)

		code := 0
		rHeader := rspHeader{}
		err = gout.New().SetRequest(req).Code(&code).BindHeader(&rHeader).Do()

		assert.NoError(t, err)
		assert.Equal(t, code, 200, fmt.Sprintf("server address:%s", router.URL))
		assert.Equal(t, rHeader.ContentEncoding, "gzip")
	}
}
