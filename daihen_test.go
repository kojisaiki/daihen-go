package main

import (
	"fmt"
	"os"
	"testing"
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

	os.Setenv("DAIHEN_RECEIVE_PORT", string(8080))
	os.Setenv("DAIHEN_PROXY_HOST", "localhost")
	os.Setenv("DAIHEN_PROXY_PORT", string(8081))
	os.Setenv("DAIHEN_PROXY_USER", "foo")
	os.Setenv("DAIHEN_PROXY_PASS", "bar")

	go daihen()
}
