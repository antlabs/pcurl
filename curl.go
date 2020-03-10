package pcurl

import (
	"github.com/guonaihong/clop"
	"github.com/guonaihong/gout"
	"net/http"
)

type Curl struct {
	Method string   `clop:"-X; --request" usage:"Specify request command to use"`
	Header []string `clop:"-H; --header" usage:"Pass custom header(s) to server"`
	Data   string   `clop:"-d; --data"   usage:"HTTP POST data"`
	URL    string   `clop:"args=url" usage:"url"`

	Err error
}

func Parse(curl string) *Curl {
	c := Curl{}
	p := clop.New([]string{} /*todo*/).SetExit(false)
	c.Err = p.Bind(&c)
	return &c
}

func ParseAndRequest(curl string) (*http.Request, error) {
	return Parse(curl).Request()
}

func (c *Curl) Request() (*http.Request, error) {
	if len(c.Method) == 0 && len(c.Data) > 0 {
		c.Method = "POST"
	}

	// TODO open
	//err := gout.New().Method(c.Method).URL(c.URL).SetBody(c.Data).Do()

	return nil, nil
}
