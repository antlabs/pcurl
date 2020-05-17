package pcurl

import (
	"testing"

	"github.com/guonaihong/gout"
	"github.com/stretchr/testify/assert"
)

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
		/*
			testHeader{
				curlHeader: []string{`curl`, "-X", "POST", `-H`, `Connection: keep-alive`, `-H`, `Cache-Control: max-age=0`, `-H`, `Upgrade-Insecure-Requests: 1`, `-H`, `User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/80.0.3987.87 Chrome/80.0.3987.87 Safari/537.36`, `-H`, `Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9`, `-H`, `Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7`, `-H`, `Cookie: username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13aeaa51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; timezone=8; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589004902,1589265192,1589341266,1589717172; Hm_lpvt_12d9f8f1740b76bb88c6691ea1672d8b=1589719525`},
				need: H{
					"Connection":                `keep-alive`,
					"Cache-Control":             `max-age=0`,
					"Upgrade-Insecure-Requests": `1`,
					"User-Agent":                `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/80.0.3987.87 Chrome/80.0.3987.87 Safari/537.36`,
					"Accept":                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9`,
					"Accept-Language":           `en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7`,
					"Cookie":                    `username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13aeaa51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; timezone=8; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589004902,1589265192,1589341266,1589717172; Hm_lpvt_12d9f8f1740b76bb88c6691ea1672d8b=1589719525`,
				},
			},
		*/
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
