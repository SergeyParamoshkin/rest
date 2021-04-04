package client

import (
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

func (u *User) isExist() (bool, error) {
	return true, nil
}

func (c *Client) Ping() (string, error) {
	req, err := http.NewRequest("GET", c.Addr+"/ping", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}
