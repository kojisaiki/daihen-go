package main

import (
	"os"
	"testing"
)

func TestDaihen(t *testing.T) {

	os.Setenv("DAIHEN_RECEIVE_PORT", string(8080))
	os.Setenv("DAIHEN_PROXY_HOST", "localhost")
	os.Setenv("DAIHEN_PROXY_PORT", string(8081))
	os.Setenv("DAIHEN_PROXY_USER", "foo")
	os.Setenv("DAIHEN_PROXY_PASS", "bar")

	go daihen()
}
