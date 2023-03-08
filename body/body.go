package body

import (
	"bytes"
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v3"
)

var ErrJSONValid = errors.New("valid json")

func valid[T string | []byte](t T) bool {
	var x any = t

	var b []byte
	switch v := x.(type) {
	case string:
		// TODO使用强制类型转换
		b = []byte(v)
	case []byte:
		b = v
	}

	b = bytes.TrimSpace(b)
	if len(b) <= 1 {
		return false
	}

	if !(b[0] == '{' && b[len(b)-1] == '}' || b[0] == '[' && b[len(b)-1] == ']') {
		return false
	}

	return json.Valid(b)
}

func jsonUnmarshal(b []byte) (o map[string]any, a []any, err error) {
	if !valid(b) {
		return nil, nil, ErrJSONValid
	}

	if err := json.Unmarshal(b, &o); err == nil /*没有错误说明是json 对象字符串*/ {
		return o, nil, nil
	}

	// 可能是array对象
	err = json.Unmarshal(b, &a)
	return nil, a, err
}

func yamlUnmarshal(bytes []byte) (o map[string]any, a []any, err error) {

	err = yaml.Unmarshal(bytes, &o)
	if err != nil {
		if err = yaml.Unmarshal(bytes, &a); err != nil {
			return nil, nil, err
		}
	}
	return
}
