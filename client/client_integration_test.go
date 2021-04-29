// client_integration_test.go
// +build integration

package client

import (
	"net/http"
	"testing"
)

var c = Client{
	Addr:   "http://localhost:3333",
	Client: http.Client{},
}

func TestPing(t *testing.T) {
	if s, err := c.Ping(); err != nil || s != "pong" {
		t.Fail()
	}
}
