# pcurl
[![Go](https://github.com/antlabs/pcurl/workflows/Go/badge.svg)](https://github.com/antlabs/pcurl/actions)
[![codecov](https://codecov.io/gh/antlabs/pcurl/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/pcurl)

pcurl是解析curl表达式的库

# feature
* 支持-X; --request，作用设置GET或POST的选项
* 支持-H; --header选项，curl中用于设置http header的选项
* 支持-d; --data选项，作用设置http body
* 支持--data-raw选项，curl用于设置http body
* 支持-F --form选项，用作设置formdata
* 支持--url选项，curl中设置url，一般不会设置这个选项
* 支持--compressed选项
* 支持-k, --insecure选项
* 支持-G, --get选项
* 支持--data-urlencode选项
* 支持内嵌到你的结构体里面，让你的cmd秒变curl

# 内容
- [一、json](#json)
- [二、form data](#form-data)
- [三、dump to json](#dump-to-json)
- [四、dump struct](#dump-struct)
# quick start
```go
package main

import (
    "fmt"
    "github.com/antlabs/pcurl"
    //"github.com/guonaihong/gout"
    "io"
    "io/ioutil"
    "net/http"
)

func main() {
    req, err := pcurl.ParseAndRequest(`curl -X POST -d 'hello world' www.qq.com`)
    if err != nil {
        fmt.Printf("err:%s\n", err)
        return
    }

    resp, err := http.DefaultClient.Do(req)
    n, err := io.Copy(ioutil.Discard, resp.Body)
    fmt.Println(err, "resp.size = ", n)

    /*
        resp := ""
        err = gout.New().SetRequest(req).BindBody(&resp).Do()

        fmt.Println(err, "resp.size = ", len(resp))
    */
}

```
## json
```go
package main

import (
    "fmt"
    "github.com/antlabs/pcurl"
    "io"
    "net/http"
    "os"
)

func main() {
    req, err := pcurl.ParseAndRequest(`curl -XPOST -d '{"hello":"world"}' 127.0.0.1:1234`)
    if err != nil {
        fmt.Printf("err:%s\n", err)
        return
    }   

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Printf("err:%s\n", err)
        return
    }   
    defer resp.Body.Close()

    io.Copy(os.Stdout, resp.Body)
}

```

## form data
```go
package main

import (
    "fmt"
    "github.com/antlabs/pcurl"
    "io"
    "net/http"
    "os"
)

func main() {
    req, err := pcurl.ParseAndRequest(`curl -XPOST -F mode=A -F text='Good morning' 127.0.0.1:1234`)
    if err != nil {
        fmt.Printf("err:%s\n", err)
        return
    }   

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Printf("err:%s\n", err)
        return
    }   
    defer resp.Body.Close()

    io.Copy(os.Stdout, resp.Body)
}

```
## dump to json
```go
package main

import (
    "fmt"
    "github.com/antlabs/pcurl"
    "io"
    "net/http"
    "os"
)

func main() {
    all, err := pcurl.ParseAndJSON(`curl https://api.openai.com/v1/completions -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_API_KEY' -d '{ "model": "text-davinci-003", "prompt": "Say this is a test", "max_tokens": 7, "temperature": 0 }'`)
	fmt.Printf("%s\n", all)
/*
{
  "url": "https://api.openai.com/v1/completions",
  "encode": {
    "body": "json"
  },
  "body": {
    "max_tokens": 7,
    "model": "text-davinci-003",
    "prompt": "Say this is a test",
    "temperature": 0
  },
  "header": [
    "Content-Type: application/json",
    "Authorization: Bearer YOUR_API_KEY"
  ]
}
}
*/
```
## dump struct 
```go
package main

import (
    "fmt"
    "github.com/antlabs/pcurl"
    "io"
    "net/http"
    "os"
)

func main() {
    all, err := pcurl.ParseAndObj(`curl https://api.openai.com/v1/completions -H 'Content-Type: application/json' -H 'Authorization: Bearer YOUR_API_KEY' -d '{ "model": "text-davinci-003", "prompt": "Say this is a test", "max_tokens": 7, "temperature": 0 }'`)

	fmt.Printf("%s\n", all)
/*
&pcurl.Req{Method:"POST", URL:"https://api.openai.com/v1/completions", Encode:pcurl.Encode{Body:"json"}, Body:map[string]interface {}{"max_tokens":7, "model":"text-davinci-003", "prompt":"Say this is a test", "temperature":0}, Header:[]string{"Content-Type: application/json", "Authorization: Bearer YOUR_API_KEY"}}
*/
```
## 继承pcurl的选项(curl)--让你的cmd秒变curl
自定义的Gen命令继续pcurl所有特性，在此基础加些自定义选项。
```go
type Gen struct {
    //curl选项
	pcurl.Curl

    //自定义选项
	Connections string        `clop:"-c; --connections" usage:"Connections to keep open"`
	Duration    time.Duration `clop:"--duration" usage:"Duration of test"`
	Thread      int           `clop:"-t; --threads" usage:"Number of threads to use"`
	Latency     string        `clop:"--latency" usage:"Print latency statistics"`
	Timeout     time.Duration `clop:"--timeout" usage:"Socket/request timeout"`
}

func main() {
	g := &Gen{}

	clop.Bind(&g)

    // pcurl包里面提供
	req, err := g.SetClopAndRequest(clop.CommandLine)
	if err != nil {
		panic(err.Error())
	}

    // 已经拿到http.Request对象
    // 如果是标准库直接通过Do()方法发送
    // 如果是裸socket，可以通过http.DumpRequestOut先转成[]byte再发送到服务端
    fmt.Printf("%p\n", req)
}

```
