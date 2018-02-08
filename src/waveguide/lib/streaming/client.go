package streaming

import (
	"encoding/json"
	"net/http"
	"waveguide/lib/config"
	"waveguide/lib/log"
)

type Client struct {
	conf *config.Config
}

func (cl *Client) Online() (streams []StreamInfo) {
	resp, err := http.Get("http://" + cl.conf.ApiServer.Addr + "/api/v1/streams/")
	if err == nil {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&streams)
	}
	if err != nil {
		log.Errorf("Failed to get online stream list: %s", err.Error())
	}
	return
}

func (cl *Client) Find(key string) (stream *StreamInfo) {
	resp, err := http.Get("http://" + cl.conf.ApiServer.Addr + "/api/v1/stream/info/" + key)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var s StreamInfo
			err = json.NewDecoder(resp.Body).Decode(&s)
			if err == nil {
				stream = &s
			}
		}
	}
	if err != nil {
		log.Errorf("Failed to find stream: %s", err.Error())
	}
	return
}

func NewClient(c *config.Config) *Client {
	return &Client{
		conf: c,
	}
}
