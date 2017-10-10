package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	url string
}

func (cl *Client) createURL(method, callbackURL string) *url.URL {
	u, _ := url.Parse(cl.url)
	u.Path = "/api/" + method
	q := u.Query()
	q.Add(ParamCallbackURL, callbackURL)
	u.RawQuery = q.Encode()
	return u
}

func (cl *Client) getHTTP() *http.Client {
	// TODO: pooling
	return new(http.Client)
}

func (cl *Client) doRequest(r *http.Request) error {
	c := cl.getHTTP()
	resp, err := c.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad backend status code: %d", resp.StatusCode)
	}
	j := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err == nil {
		emsg, ok := j["Error"]
		if ok {
			if emsg != nil {
				err = fmt.Errorf("%s", emsg)
			}
		}
	}
	return err
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}
