package pcurl

import (
	"fmt"
	"github.com/guonaihong/clop"
	"github.com/guonaihong/gout"
	"net/http"
	"os"
	"strings"
)

type Curl struct {
	Method string   `clop:"-X; --request" usage:"Specify request command to use"`
	Header []string `clop:"-H; --header" usage:"Pass custom header(s) to server"`
	Data   string   `clop:"-d; --data"   usage:"HTTP POST data"`
	Form   string   `clop:"-F; --form" usage:"Specify multipart MIME data"`
	URL    string   `clop:"args=url" usage:"url"`

	Err error
}

func ParseAndRequest(curl string) (*http.Request, error) {
	return Parse(curl).Request()
}

func Parse(curl string) *Curl {
	return ParseSlice([]string{})
}

func ParseSlice(curl []string) *Curl {
	c := Curl{}
	if len(curl) > 0 && strings.ToLower(curl[0]) == "curl" {
		curl = curl[1:]
	}
	fmt.Printf("%#v\n", curl)

	p := clop.New(curl).SetExit(false)
	c.Err = p.Bind(&c)
	return &c
}

func (c *Curl) createHeader() []string {
	if len(c.Header) == 0 {
		return nil
	}

	header := make([]string, len(c.Header)*2)
	index := 0
	for _, v := range header {
		hv := strings.Split(v, ":")
		if len(hv) != 2 {
			continue
		}

		header[index] = hv[0]
		index++
		header[index] = hv[1]
		index++
	}

	return header
}

func (c *Curl) Request() (*http.Request, error) {
	if len(c.Method) == 0 && len(c.Data) > 0 {
		c.Method = "POST"
	}

	var (
		data interface{}
	)

	header := c.createHeader()

	data = c.Data
	if len(c.Data) > 0 && c.Data[0] == '@' {
		fd, err := os.Open(c.Data[1:])
		if err != nil {
			return nil, err
		}

		defer fd.Close()

		data = fd
	}

	g := gout.New().
		SetMethod(c.Method). //设置method POST or GET or DELETE
		Debug(true)          //打开debug模式

	if header != nil {
		g.SetHeader(header) //设置http header
	}

	return g.SetURL(c.URL). //设置url
				SetBody(data). //设置http body
				Request()      //获取*http.Request
}
