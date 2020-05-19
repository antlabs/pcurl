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
