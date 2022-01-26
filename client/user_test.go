// client_test.go
// +build !integration

package client_test

import (
	"testing"

	"github.com/SergeyParamoshkin/rest/client"
)

func TestUserIsExist(t *testing.T) {
	u := client.User{
		ID:   0,
		Name: "test",
	}

	t.Parallel()

	if ok, err := u.IsExist(); !ok || err != nil {
		t.Fail()
	}
}
