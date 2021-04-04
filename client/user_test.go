// client_test.go
// +build !integration

package client

import "testing"

func TestUserIsExist(t *testing.T) {
	u := User{}
	if ok, err := u.isExist(); !ok || err != nil {
		t.Fail()
	}
}
