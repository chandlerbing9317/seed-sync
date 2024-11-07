package common

import (
	"net/http"
	"net/url"
)

func GetProxyFunc(proxyUrl string) func(req *http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl)
	}
}
