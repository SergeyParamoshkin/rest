package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	http.Client
	Addr string
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (u *User) IsExist() (bool, error) {
	return true, nil
}

func (c *Client) Ping() (string, error) {
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Addr+"/ping", nil)
	if err != nil {
		return "", fmt.Errorf("%w: newRequest", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: response", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: readall", err)
	}

	return string(body), err
}
