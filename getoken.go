package pcurl

import (
	"strings"
)

const (
	Unused       = iota
	DoubleQuotes //双引号
	SingleQuotes //单引号
	WordEnd
)

// TODO 考虑转义
// TODO 各种换行符号
func GetArgsToken(curlString string) (curl []string) {
	var (
		sign      int
		needSpace bool
		word      int
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

	return
}
