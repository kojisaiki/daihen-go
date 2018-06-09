package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

const (
	ProxyAuthHeader = "Proxy-Authorization"
)

func SetBasicAuth(username, password string, req *http.Request) {
	req.Header.Set(ProxyAuthHeader, fmt.Sprintf("Basic %s", basicAuth(username, password)))
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func main() {
	username, password := "foo", "bar"
	proxyhost, proxyport := "localhost", "8080"

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	connectReqHandler := func(req *http.Request) {
		SetBasicAuth(username, password, req)
	}
	proxy.ConnectDial = proxy.NewConnectDialToProxyWithHandler(fmt.Sprintf("http://%s:%s", proxyhost, proxyport), connectReqHandler)
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		fmt.Println("---------------request!")
		SetBasicAuth(username, password, r)
		return r, nil
	})

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
