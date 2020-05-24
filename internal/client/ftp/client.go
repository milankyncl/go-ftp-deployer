package ftp

import (
	"github.com/secsy/goftp"
	"os"
	"time"
)

type Client struct {
	client *goftp.Client
}

func NewClient(host string, user string, password string) (*Client, error) {
	config := goftp.Config{
		User:               user,
		Password:           password,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
	}
	c, err := goftp.DialConfig(config, host)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: c,
	}, nil
}

func (c *Client) Upload(locPath string, extPath string) error {
	file, err := os.Open(locPath)
	if err != nil {
		return err
	}
	return c.client.Store(extPath, file)
}

func (c *Client) CreateDir(path string) error {
	_, err := c.client.Mkdir(path)
	return err
}

func (c *Client) FileExists(path string) bool {
	_, err := c.client.Stat(path)
	return err == nil
}

func (c *Client) Close() {
	err := c.client.Close()
	if err != nil {
		panic(err)
	}
}
