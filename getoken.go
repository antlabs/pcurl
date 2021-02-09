// Copyright [2020-2021] [guonaihong]
// Apache 2.0
package pcurl

import (
	"errors"
	"strings"
)

var (
	ErrSingleQuotes = errors.New("unquoted single quote")
	ErrDoubleQuotes = errors.New("unquoted double quote")
	ErrUnknown      = errors.New("pcurl:GetArgsToken:unknown error")
)

type Sign int

const (
	Unused       Sign = iota
	DoubleQuotes      //双引号
	SingleQuotes      //单引号
	WordEnd
)

// TODO 考虑转义
// TODO 各种换行符号
func GetArgsToken(curlString string) (curl []string, err error) {
	var (
		sign      Sign
		needSpace bool
		word      Sign
	)

	var buf strings.Builder

	for _, b := range curlString {

		// 先处理有作用域的符号
		if sign == Unused {
			switch b {
			case '"':
				needSpace = true
				sign = DoubleQuotes
				continue
			case '\'':
				needSpace = true
				sign = SingleQuotes
				continue
			}
		}

		// 处理作用域右边的"
		if b == '"' {
			if sign == DoubleQuotes {
				sign = Unused
				needSpace = false
				continue
			}
		}

		// 处理作用域右边的'
		if b == '\'' {
			if sign == SingleQuotes {
				sign = Unused
				needSpace = false
				continue
			}
		}

		//
		if !needSpace && b == ' ' {
			word = WordEnd
			continue
		}

		//fmt.Printf("(%c):%d, %t\n", b, sign, word == WordEnd)

		if word == WordEnd {
			curl = append(curl, buf.String())
			buf.Reset()
			//sign = Unused
			word = Unused
			//needSpace = false
		}

		buf.WriteRune(b)

	}

	if buf.Len() > 0 {
		curl = append(curl, buf.String())
		buf.Reset()
	}

	if sign != Unused {
		return nil, toErr(sign)
	}

	return
}

func toErr(sign Sign) error {
	switch sign {
	case SingleQuotes:
		return ErrSingleQuotes
	case DoubleQuotes:
		return ErrDoubleQuotes
	case Unused:
		return nil
	default:
		return ErrUnknown
	}
}
