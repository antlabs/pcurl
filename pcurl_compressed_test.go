package pcurl

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
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
		if err != nil {
			t.Fatalf("ParseAndRequest failed: %v", err)
		}

		code := 0
		rHeader := rspHeader{}
		err = gout.New().SetRequest(req).Code(&code).BindHeader(&rHeader).Do()
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if code != 200 {
			t.Fatalf("unexpected status code, got=%d want=%d server=%s", code, 200, router.URL)
		}
		if rHeader.ContentEncoding != "gzip" {
			t.Fatalf("unexpected Content-Encoding, got=%q want=%q", rHeader.ContentEncoding, "gzip")
		}
	}
}
