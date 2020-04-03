# pcurl
pcurl是解析curl表达式的库，还在继续完善中。。。

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