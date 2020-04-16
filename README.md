# pcurl
[![Go](https://github.com/antlabs/pcurl/workflows/Go/badge.svg)](https://github.com/antlabs/pcurl/actions)
[![codecov](https://codecov.io/gh/guonaihong/pcurl/branch/master/graph/badge.svg)](https://codecov.io/gh/guonaihong/pcurl)

pcurl是解析curl表达式的库，还在继续完善中。。。

# feature
* 支持-X; --request，作用设置GET或POST的选项
* 支持-H; --header选项，curl中用于设置http header的选项
* 支持-d; --data选项，作用设置http body
* 支持--data-raw选项，curl用于设置http body
* 支持-F --form选项，用作设置formdata
* 支持--url选项，curl中设置url，一般不会设置这个选项

# 内容
- [json](#json)
- [form data](#form-data)
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
