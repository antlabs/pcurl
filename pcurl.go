// Copyright [2020-2021] [guonaihong]
// Apache 2.0
package pcurl

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antlabs/pcurl/body"
	"github.com/guonaihong/clop"
	"github.com/guonaihong/gout"
)

// Curl结构体
type Curl struct {
	Method        string   `clop:"-X; --request" usage:"Specify request command to use"`
	Get           bool     `clop:"-G; --get" usage:"Put the post data in the URL and use GET"`
	Header        []string `clop:"-H; --header" usage:"Pass custom header(s) to server"`
	Data          string   `clop:"-d; --data"   usage:"HTTP POST data"`
	DataRaw       string   `clop:"--data-raw" usage:"HTTP POST data, '@' allowed"`
	Form          []string `clop:"-F; --form" usage:"Specify multipart MIME data"`
	URL2          string   `clop:"args=url2" usage:"url2"`
	URL           string   `clop:"--url" usage:"URL to work with"`
	Location      bool     `clop:"-L; --location" usage:"Follow redirects"` //TODO
	DataUrlencode []string `clop:"--data-urlencode" usage:"HTTP POST data url encoded"`

	Compressed bool `clop:"--compressed" usage:"Request compressed response"`
	// 在响应包里面打印http header, 仅做字段赋值
	Include  bool `clop:"-i;--include" usage:"Include the HTTP response headers in the output. The HTTP response headers can include things like server name, cookies, date of the document, HTTP version and more."`
	Insecure bool `clop:"-k; --insecure" usage:"Allow insecure server connections when using SSL"`
	Err      error
	p        *clop.Clop
}

const (
	bodyEmpby     = "empty"
	bodyURLEncode = "data-urlencode"
	bodyForm      = "form"
	bodyData      = "data"
	bodyDataRaw   = "data-raw"
)

// 内嵌Curl结构体至cmd里面使用，让你的cmd秒变为curl
func (c *Curl) SetClopAndRequest(clop2 *clop.Clop) (*http.Request, error) {
	if c == nil {
		return nil, errors.New("curl is nil")
	}

	c.p = clop2

	return c.Request()
}

// ParseString是链式API结构, 如果要拿*http.Request，后接Request()即可
func ParseString(curl string) *Curl {
	c := Curl{}
	curlSlice, err := GetArgsToken(curl)
	c.Err = err
	return parseSlice(curlSlice, &c)
}

// ParseSlice和ParseString的区别，ParseSlice里面保存解析好的curl表达式
func ParseSlice(curl []string) *Curl {
	c := Curl{}
	return parseSlice(curl, &c)
}

func (c *Curl) createHeader() []string {
	if len(c.Header) == 0 {
		return nil
	}

	header := make([]string, len(c.Header)*2)
	index := 0
	for _, v := range c.Header {
		pos := strings.IndexByte(v, ':')
		if pos == -1 {
			continue
		}

		header[index] = v[:pos]
		index++
		header[index] = v[pos+1:]
		index++
	}

	return header
}

func (c *Curl) getBodyEncodeAndObj() (string, any, error) {
	name := c.findHighestPriority()
	if name == bodyURLEncode {
		allQuery := strings.Join(c.DataUrlencode, "&")
		u, err := url.ParseQuery(allQuery)
		return body.EncodeWWWForm, u, err
	}

	if name == bodyData {
		return body.Unmarshal([]byte(c.Data))
	}

	if name == bodyEmpby {
		return "", nil, nil
	}

	return "", nil, ErrUnknownEncode
}

func (c *Curl) findHighestPriority() string {

	// 获取 --data-urlencoded,-F or --form, -d or --data, --data-raw的命令行优先级别
	// 绑定body和优先级的关系
	bodyName := bodyEmpby
	max := uint64(0)
	if index, ok := c.p.GetIndexEx(bodyURLEncode); ok && index > max {
		bodyName = bodyURLEncode
	}

	if index, ok := c.p.GetIndexEx(bodyForm); ok && index > max {
		bodyName = bodyForm
	}

	if index, ok := c.p.GetIndexEx(bodyData); ok && index > max {
		bodyName = bodyData
	}

	if index, ok := c.p.GetIndexEx(bodyDataRaw); ok && index > max {
		bodyName = bodyDataRaw
	}

	return bodyName
}

func (c *Curl) createWWWForm() ([]any, error) {
	if len(c.DataUrlencode) == 0 {
		return nil, nil
	}

	form := make([]any, len(c.DataUrlencode)*2)
	index := 0
	for _, v := range c.DataUrlencode {
		pos := strings.IndexByte(v, '=')
		if pos == -1 {
			continue
		}

		form[index] = v[:pos]
		index++
		form[index] = v[pos+1:]
		index++
	}

	return form, nil
}

func (c *Curl) createForm() ([]any, error) {
	if len(c.Form) == 0 {
		return nil, nil
	}

	form := make([]any, len(c.Form)*2)
	index := 0
	for _, v := range c.Form {
		pos := strings.IndexByte(v, '=')
		if pos == -1 {
			continue
		}

		form[index] = v[:pos]
		index++
		fieldValue := v[pos+1:]
		if len(fieldValue) > 0 && fieldValue[0] == '@' {

			form[index] = gout.FormFile(fieldValue[1:])
		} else {

			form[index] = fieldValue
		}

		index++
	}

	return form, nil
}

func (c *Curl) getURL() string {
	url := c.URL2
	if c.p.GetIndex("url") > c.p.GetIndex("url2") {
		url = c.URL
	}

	return url
}

func (c *Curl) setMethod() {
	// 在curl里面-X的选项的优先级别比-G高，所以c.Method为空时才会看c.Get是否设置
	if len(c.Method) == 0 && c.Get {
		c.Method = "GET"
		return
	}

	if len(c.Method) != 0 {
		return
	}

	if len(c.Data) > 0 {
		c.Method = "POST"
		return
	}

	c.Method = "GET"
}

func (c *Curl) Request() (req *http.Request, err error) {

	var (
		data    any
		form    []any
		wwwForm []any
		dataRaw string
	)

	defer func() {
		if c.Err != nil {
			err = c.Err
		}
	}()

	c.setMethod()

	header := c.createHeader()

	switch c.findHighestPriority() {
	case bodyURLEncode:
		if wwwForm, err = c.createWWWForm(); err != nil {
			return nil, err
		}
	case bodyForm:
		if form, err = c.createForm(); err != nil {
			return nil, err
		}
	case bodyData:
		dataRaw = c.Data
	case bodyDataRaw:
		dataRaw = c.DataRaw
	}

	var hc *http.Client

	if c.Insecure {
		hc = &defaultInsecureSkipVerify
	}

	client := gout.New(hc)
	url := c.getURL()

	// TODO gout SetURL放在最后一项设置
	client.SetURL(url) //设置url

	client.SetMethod(c.Method) //设置method POST or GET or DELETE

	if c.Compressed {
		header = append(header, "Accept-Encoding", "deflate, gzip")
	}

	if len(dataRaw) > 0 {
		data = dataRaw
	}
	if len(dataRaw) > 0 && dataRaw[0] == '@' {
		fd, err := os.Open(dataRaw[1:])
		if err != nil {
			return nil, err
		}

		defer fd.Close()

		data = fd
	}

	if len(header) > 0 {
		client.SetHeader(header) //设置http header
	}

	if len(form) > 0 {
		client.SetForm(form) //设置formdata
	}

	if len(wwwForm) > 0 {
		client.SetWWWForm(wwwForm) // 设置x-www-form-urlencoded格式数据
	}

	if data != nil {
		client.SetBody(data)
	}

	return client.Request() //获取*http.Request
}

func parseSlice(curl []string, c *Curl) *Curl {
	if len(curl) > 0 && strings.ToLower(curl[0]) == "curl" {
		curl = curl[1:]
	}

	c.p = clop.New(curl).SetExit(false)
	c.Err = c.p.Bind(&c)
	return c
}
