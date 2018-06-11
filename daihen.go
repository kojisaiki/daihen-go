package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

type DaihenConfiguration struct {
	receivePort uint16
	proxyHost   string
	proxyPort   uint16
	username    string
	password    string
}

func prepareEnv() (DaihenConfiguration, bool) {
	receivePort := os.Getenv("DAIHEN_RECEIVE_PORT")
	proxyHost := os.Getenv("DAIHEN_PROXY_HOST")
	proxyPort := os.Getenv("DAIHEN_PROXY_PORT")
	username := os.Getenv("DAIHEN_PROXY_USER")
	password := os.Getenv("DAIHEN_PROXY_PASS")

	if receivePort == "" || proxyHost == "" || proxyPort == "" || username == "" || password == "" {
		fmt.Println(`Require following environment variables:
  - DAIHEN_RECEIVE_PORT:Daihen listen this port.
  - DAIHEN_PROXY_HOST:Daihen bypass request to this hostname.
  - DAIHEN_PROXY_PORT:Daihen bypass request to this port.
  - DAIHEN_PROXY_USER:Basic authentication username for proxy.
  - DAIHEN_PROXY_PASS:Basic authentication password for proxy.`)
		return DaihenConfiguration{}, false
	}

	config := DaihenConfiguration{}
	_receivePort, _ := strconv.ParseUint(receivePort, 10, 16)
	config.receivePort = uint16(_receivePort)
	config.proxyHost = proxyHost
	_proxyPort, _ := strconv.ParseUint(proxyPort, 10, 16)
	config.proxyPort = uint16(_proxyPort)
	config.username = username
	config.password = password
	return config, true
}

func daihen() {

	config, ok := prepareEnv()

	if !ok {
		return
	}

	proxy := goproxy.NewProxyHttpServer()

	connectReqHandler := func(req *http.Request) {
		SetBasicAuth(config.username, config.password, req)
	}
	proxy.ConnectDial = proxy.NewConnectDialToProxyWithHandler(fmt.Sprintf("http://%s:%d", config.proxyHost, config.proxyPort), connectReqHandler)
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		SetBasicAuth(config.username, config.password, r)
		return r, nil
	})

	fmt.Printf("serve on:%d¥n", config.receivePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.receivePort), proxy))
}

func main() {
	daihen()
}
