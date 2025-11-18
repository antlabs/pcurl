package pcurl

import (
	"testing"
)

func Test_Issue9(t *testing.T) {

	req, err := ParseAndRequest(`curl 'http://xxxx.cc/admin/index.php?act=admin' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' -H 'Upgrade-Insecure-Requests: 1' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/80.0.3987.87 Chrome/80.0.3987.87 Safari/537.36' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' -H 'Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7' -H 'Cookie: username=admin; token=b7ea3ec643e4ea4871dfe515c559d28bc0d23b6d9d6b22daf206f1de9aff13e51591323199; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D; timezone=8; Hm_lvt_12d9f8f1740b76bb88c6691ea1672d8b=1589004902,1589265192,1589341266,1589717172; Hm_lpvt_12d9f8f1740b76bb88c6691ea1672d8b=1589719525'`)
	if err != nil {
		t.Fatalf("ParseAndRequest failed: %v", err)
	}
	if req == nil {
		t.Fatalf("expected non-nil request")
	}
	if len(req.Header) == 0 {
		t.Fatalf("expected non-empty headers")
	}
}
