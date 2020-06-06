package ftp

import (
	"github.com/secsy/goftp"
	"io"
	"os"
	"path"
	"time"
)

type Client struct {
	client  *goftp.Client
	rootDir string
}

func NewClient(host string, user string, password string, rootDir string) (*Client, error) {
	config := goftp.Config{
		User:               user,
		Password:           password,
		ConnectionsPerHost: 15,
		Timeout:            10 * time.Second,
	}
	c, err := goftp.DialConfig(config, host)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:  c,
		rootDir: rootDir,
	}, nil
}

func (c *Client) Upload(locPath string, extPath string) error {
	file, err := os.Open(locPath)
	if err != nil {
		return err
	}
	return c.client.Store(
		path.Join(c.rootDir, extPath),
		file,
	)
}

func (c *Client) CreateDir(dirname string) error {
	_, err := c.client.Mkdir(
		path.Join(c.rootDir, dirname),
	)
	return err
}

func (c *Client) FileExists(filename string) bool {
	_, err := c.client.Stat(
		path.Join(c.rootDir, filename),
	)
	return err == nil
}

func (c *Client) ReadFile(filename string) (*File, error) {
	dest := File{}
	err := c.client.Retrieve(
		path.Join(c.rootDir, filename),
		&dest,
	)
	return &dest, err
}

func (c *Client) WriteFile(filename string, src io.Reader) error {
	return c.client.Store(
		path.Join(c.rootDir, filename),
		src,
	)
}

func (c *Client) Close() {
	err := c.client.Close()
	if err != nil {
		panic(err)
	}
}
