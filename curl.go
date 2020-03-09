package pcurl

import (
	"github.com/guonaihong/clop"
	"net/http"
)

type Curl struct {
	Method string   `clop:"-X; --request" usage:"Specify request command to use"`
	Header []string `clop:"-H; --header" usage:"Pass custom header(s) to server"`
	Data   string   `clop:"-d; --data"   usage:"HTTP POST data"`

	Err error
}

func Parse(curl string) *Curl {
	c := Curl{}
	p := clop.New([]string{} /*todo*/).SetExit(false)
	c.Err = p.Bind(&c)
	return &c
}

func (c *Curl) Request() (*http.Request, error) {
	return nil, nil
}
