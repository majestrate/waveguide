package streaming

import (
	"waveguide/lib/config"
)

type Client struct {
	conf *config.Config
}

func (cl *Client) Online() (streams []string) {
	return
}

func (cl *Client) Find(key string) (stream *StreamInfo) {
	return
}

func NewClient(c *config.Config) *Client {
	return &Client{
		conf: c,
	}
}
