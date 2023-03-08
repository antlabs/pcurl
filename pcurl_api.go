package pcurl

import (
	"encoding/json"
	"net/http"
)

// 解析curl字符串形式表达式，并返回*http.Request
func ParseAndRequest(curl string) (*http.Request, error) {
	return ParseString(curl).Request()
}

// 解析curl字符串形式表达式，并返回结构体
func ParseAndObj(curl string) (r *Req, err error) {

	r = &Req{}
	c := ParseString(curl)
	if c.Err != nil {
		return nil, err
	}
	c.setMethod()

	r.Method = c.Method
	r.URL = c.URL2
	r.Header = c.Header

	r.Encode.Body, r.Body, err = c.getBodyEncodeAndObj()
	return
}

func ParseAndJSON(curl string) (jsonBytes []byte, err error) {
	o, err := ParseAndObj(curl)
	if err != nil {
		return nil, err
	}

	return json.Marshal(o)
}
