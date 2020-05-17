package pcurl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试错误情况
func Test_GetArgsTokenFail(t *testing.T) {
	type TestArgs struct {
		in string
	}

	for _, v := range []TestArgs{
		TestArgs{
			in: `'hello`,
		},
		TestArgs{
			in: `"hello`,
		},
	} {

		_, err := GetArgsToken(v.in)
		assert.Error(t, err)
	}
}

// 测试正确的情况
func Test_GetArgsToken(t *testing.T) {
	type TestArgs struct {
		in   string
		need []string
	}

	for _, v := range []TestArgs{
		TestArgs{
			in:   `curl -XGET "http://192.168.6.100:9200/eval-log/_search" -H 'Content-Type: application/json' -d'{  "query": {    "match": {      "level": "error"    }  }}'`,
			need: []string{`curl`, `-XGET`, "http://192.168.6.100:9200/eval-log/_search", "-H", `Content-Type: application/json`, `-d{  "query": {    "match": {      "level": "error"    }  }}`},
		},
		TestArgs{
			in:   `curl --location --request DELETE '192.168.5.213:10010/delete/rule?appkey=xx' --header 'Content-Type: text/plain' --data-raw '{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416}'`,
			need: []string{`curl`, `--location`, `--request`, `DELETE`, `192.168.5.213:10010/delete/rule?appkey=xx`, `--header`, `Content-Type: text/plain`, `--data-raw`, `{"type":"region","region":"bj","business":"asr","protocol":"private","connect":416}`},
		},
		TestArgs{
			in:   `'{"s":"{\"s\":\"S\"}"}'`,
			need: []string{`{"s":"{\"s\":\"S\"}"}`},
		},
		TestArgs{
			in:   `curl 'http://xxx.cc/admin/index.php?act=admin' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' -H 'Upgrade-Insecure-Requests: 1' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/80.0.3987.87 Chrome/80.0.3987.87 Safari/537.36' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9' -H 'Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7' -H 'Cookie: username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13aeaa51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; timezone=8; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589004902,1589265192,1589341266,1589717172; Hm_lpvt_12d9f8f1740b76bb88c6691ea1672d8b=1589719525'`,
			need: []string{`curl`, `http://xxx.cc/admin/index.php?act=admin`, `-H`, `Connection: keep-alive`, `-H`, `Cache-Control: max-age=0`, `-H`, `Upgrade-Insecure-Requests: 1`, `-H`, `User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/80.0.3987.87 Chrome/80.0.3987.87 Safari/537.36`, `-H`, `Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9`, `-H`, `Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7`, `-H`, `Cookie: username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13aeaa51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; timezone=8; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589004902,1589265192,1589341266,1589717172; Hm_lpvt_12d9f8f1740b76bb88c6691ea1672d8b=1589719525`},
		},
	} {

		got, err := GetArgsToken(v.in)
		assert.NoError(t, err)
		assert.Equal(t, v.need, got)
	}
}
