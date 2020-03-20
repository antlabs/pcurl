package pcurl

import (
	"fmt"
	"github.com/guonaihong/clop"
	"github.com/guonaihong/gout"
	"net/http"
	"os"
)

type Curl struct {
	Method string   `clop:"-X; --request" usage:"Specify request command to use"`
	Header []string `clop:"-H; --header" usage:"Pass custom header(s) to server"`
	Data   string   `clop:"-d; --data"   usage:"HTTP POST data"`
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
	if len(curl) > 0 && curl[0] == "curl" {
		curl = curl[1:]
	}

	p := clop.New(curl).SetExit(false)
	c.Err = p.Bind(&c)
	return &c
}

func (c *Curl) Request() (*http.Request, error) {
	if len(c.Method) == 0 && len(c.Data) > 0 {
		c.Method = "POST"
	}

	var data interface{}

	data = c.Data
	if len(c.Data) > 0 && c.Data[0] == '@' {
		fd, err := os.Open(c.Data[1:])
		if err != nil {
			return nil, err
		}

		defer fd.Close()

		data = fd
	}

	return gout.New().SetMethod(c.Method).Debug(true).SetURL(c.URL).SetBody(data).Request()
}
