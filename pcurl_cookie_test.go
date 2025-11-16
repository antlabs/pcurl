package pcurl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
)

// curl -b / --cookie
func Test_Cookie_Option(t *testing.T) {
	type testCase struct {
		name      string
		curlSlice []string
		need      H
	}

	// 准备一个临时 cookie 文件（Netscape 格式）
	dir := t.TempDir()
	cookieFile := filepath.Join(dir, "cookie.txt")
	content := "# Netscape HTTP Cookie File\n" +
		"example.com\tTRUE\t/\tFALSE\t0\tfoo\tbar\n" +
		"example.com\tTRUE\t/\tFALSE\t0\tbaz\tqux\n"
	assert.NoError(t, os.WriteFile(cookieFile, []byte(content), 0o644))

	for _, c := range []testCase{
		{
			name:      "single name=value",
			curlSlice: []string{"curl", "-X", "GET", "-b", "a=1"},
			need: H{
				"Cookie": "a=1",
			},
		},
		{
			name:      "multi -b name=value",
			curlSlice: []string{"curl", "-X", "GET", "-b", "a=1", "-b", "b=2"},
			need: H{
				"Cookie": "a=1; b=2",
			},
		},
		{
			name:      "cookie file",
			curlSlice: []string{"curl", "-X", "GET", "-b", cookieFile},
			need: H{
				"Cookie": "foo=bar; baz=qux",
			},
		},
		{
			name:      "cookie file and name=value",
			curlSlice: []string{"curl", "-X", "GET", "-b", cookieFile, "-b", "c=3"},
			need: H{
				"Cookie": "foo=bar; baz=qux; c=3",
			},
		},
		{
			name:      "-H Cookie override -b",
			curlSlice: []string{"curl", "-X", "GET", "-b", "a=1", "-H", "Cookie: x=y"},
			need: H{
				"Cookie": "x=y",
			},
		},
	} {

		// 创建测试服务端
		ts := createGeneralHeader(c.need, t)

		// 解析curl表达式
		req, err := ParseSlice(append(c.curlSlice, ts.URL)).Request()
		assert.NoError(t, err, c.name)

		var got H
		code := 0
		// 发送请求
		err = gout.New().SetRequest(req).Debug(true).Code(&code).BindJSON(&got).Do()
		assert.NoError(t, err, c.name)
		assert.Equal(t, 200, code, c.name)
		assert.Equal(t, c.need, got, c.name)
	}
}
