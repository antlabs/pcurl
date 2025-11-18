package pcurl

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/clop"
	"github.com/guonaihong/gout"
)

type H map[string]any

// mapsEqual compares two H maps for equality in tests.
func mapsEqual(a, b H) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

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

func createGeneral2() *httptest.Server {
	router := func() *gin.Engine {
		router := gin.New()

		router.Any("/*test", func(c *gin.Context) {
			io.Copy(c.Writer, c.Request.Body)
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
		if err != nil {
			t.Fatalf("ParseAndRequest failed: %v", err)
		}

		got := ""
		err = gout.New().SetRequest(req).BindBody(&got).Do()
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if got != need[index] {
			t.Fatalf("unexpected method, got=%q want=%q", got, need[index])
		}
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
		if err != nil {
			t.Fatalf("ParseSlice.Request failed: %v", err)
		}

		s := ""
		//发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindBody(&s).Do()
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		if code != 200 {
			t.Fatalf("unexpected status code, got=%d want=%d", code, 200)
		}
		if s != urlData.need {
			t.Fatalf("unexpected body, got=%q want=%q", s, urlData.need)
		}
	}
}

// 测试ParseSliceAndRequest 正确的情况
func Test_ParseSliceAndRequest(t *testing.T) {
	type testParseSlice struct {
		curl []string
		need string
	}

	//在这里加更多测试数据，for + cast table，很方便加测试数据
	for _, d := range []testParseSlice{
		{
			curl: []string{"-H", "hello:word", "-H", "abc:def", "-d", "body content", "www.qq.com"},
			need: "POST / HTTP/1.1\r\n" +
				"Host: www.qq.com\r\n" +
				"User-Agent: Go-http-client/1.1\r\n" +
				"Content-Length: 12\r\n" +
				"Abc: def\r\n" +
				"Hello: word\r\n" +
				"Accept-Encoding: gzip\r\n\r\n" +
				"body content",
		},
	} {
		//声明解析器
		clop2 := clop.New(d.curl).SetExit(false)
		//声明存放解析之后的结构体
		c := Curl{}
		//解析
		clop2.Bind(&c)

		//生成req
		req, err := c.SetClopAndRequest(clop2)
		if err != nil {
			t.Fatalf("SetClopAndRequest failed: %v", err)
		}

		//把req转成[]byte
		all, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			t.Fatalf("DumpRequestOut failed: %v", err)
		}

		//比较数据看下对错
		if !reflect.DeepEqual(d.need, string(all)) {
			t.Fatalf("unexpected request dump, got=%q want=%q", string(all), d.need)
		}
	}

}

// 测试ParseSliceAndRequest 错误的情况
func Test_ParseSliceAndRequest_Error(t *testing.T) {
	c := (*Curl)(nil)
	_, err := c.SetClopAndRequest(clop.New([]string{}).SetExit(false))
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

}
