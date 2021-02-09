// Copyright [2020-2021] [guonaihong]
// Apache 2.0
package pcurl

import (
	"crypto/tls"
	"net/http"
)

var defaultInsecureSkipVerify = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}
