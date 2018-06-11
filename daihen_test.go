package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/elazarl/goproxy"
)

func TestPrepareEnv(t *testing.T) {
	_, ok := prepareEnv()
	if ok {
		fmt.Println("The case where all environments were lost was allowed.")
		t.Fail()
	}

	os.Setenv("DAIHEN_RECEIVE_PORT", "8080")
	_, ok = prepareEnv()
	if ok {
		fmt.Println("The case where 1 environments was set was allowed.")
		t.Fail()
	}

	os.Setenv("DAIHEN_PROXY_HOST", "localhost")
	_, ok = prepareEnv()
	if ok {
		fmt.Println("The case where 2 environments was set was allowed.")
		t.Fail()
	}

	os.Setenv("DAIHEN_PROXY_PORT", "8081")
	_, ok = prepareEnv()
	if ok {
		fmt.Println("The case where 3 environments was set was allowed.")
		t.Fail()
	}

	os.Setenv("DAIHEN_PROXY_USER", "foo")
	_, ok = prepareEnv()
	if ok {
		fmt.Println("The case where 4 environments was set was allowed.")
		t.Fail()
	}

	os.Setenv("DAIHEN_PROXY_PASS", "bar")
	var config DaihenConfiguration
	config, ok = prepareEnv()
	if !ok {
		fmt.Println("The case where all environments were set was not allowed.")
		t.Fail()
	}

	// expect config value
	if config.receivePort != 8080 {
		fmt.Printf("Config 'receivePort' is not %d. Actual %d.", 8080, config.receivePort)
		t.Fail()
	}
	if config.proxyHost != "localhost" {
		fmt.Printf("Config 'proxyHost' is not %s. Actual %s.", "localhost", config.proxyHost)
		t.Fail()
	}
	if config.proxyPort != 8081 {
		fmt.Printf("Config 'proxyPort' is not %d. Actual %d.", 8081, config.proxyPort)
		t.Fail()
	}
	if config.username != "foo" {
		fmt.Printf("Config 'username' is not %s. Actual %s.", "foo", config.username)
		t.Fail()
	}
	if config.password != "bar" {
		fmt.Printf("Config 'password' is not %s. Actual %s.", "bar", config.password)
		t.Fail()
	}
}

func TestDaihen(t *testing.T) {

	// setup dummy end proxy
	proxymock := goproxy.NewProxyHttpServer()
	proxymock.Verbose = true
	go http.ListenAndServe(":8081", proxymock)

	// setup subject
	os.Setenv("DAIHEN_RECEIVE_PORT", "8080")
	os.Setenv("DAIHEN_PROXY_HOST", "localhost")
	os.Setenv("DAIHEN_PROXY_PORT", "8081")
	os.Setenv("DAIHEN_PROXY_USER", "foo")
	os.Setenv("DAIHEN_PROXY_PASS", "bar")
	go daihen()

	time.Sleep(1 * time.Second)

	// fire a http request: client --> subject proxy(daihen) --> end proxy --> internet
	proxyUrl, _ := url.Parse("http://localhost:8080")
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	rsp, err := client.Get("https://golang.org/doc/tos.html")
	if err != nil {
		log.Fatalf("get rsp failed:%v", err)
		t.Fail()
	}
	defer rsp.Body.Close()
	data, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != http.StatusOK {
		log.Fatalf("status %d, data %s", rsp.StatusCode, data)
		t.Fail()
	}
	log.Printf("rsp:%s", data)
}
