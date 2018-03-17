package adn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"waveguide/lib/util"
)

type ChannelMembers struct {
	//Users     []interface{} `json:"user_ids"`
	AnyUser   bool `json:"any_user"`
	Public    bool `json:"public"`
	Immutable bool `json:"immutable"`
}

type ChannelStats struct {
	Messages    int64 `json:"messages"`
	Subscribers int64 `json:"subscribers"`
}

type Channel struct {
	Counts  ChannelStats   `json:"counts"`
	Type    string         `json:"type"`
	Readers ChannelMembers `json:"readers"`
	Writers ChannelMembers `json:"writers"`
	ID      ChanID         `json:"id"`
	//Owner   User           `json:"owner"`
}

type ChannelResponse struct {
	Data Channel      `json:"data"`
	Meta MetaResponse `json:"meta"`
}

func (c *Client) DeleteChannel(token string, chnlID ChanID) (err error) {
	u, _ := url.Parse(c.conf.Provider + fmt.Sprintf("stream/0/channels/%s", chnlID.String()))
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	if err == nil {
		var resp *http.Response
		var req *http.Request
		req, err = http.NewRequest("DELETE", u.String(), nil)
		if err == nil {
			req.Header.Set("Content-Type", "application/json; encoding=UTF-8")
			resp, err = c.http.Do(req)
			if err == nil {
				// TODO: logging
				resp.Body.Close()
			}
		}
	}
	return
}

func (c *Client) CreateChannel(token string, ch Channel) (chnl *Channel, err error) {
	u, _ := url.Parse(c.conf.Provider + "stream/0/channels")
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	buff := new(util.Buffer)
	err = json.NewEncoder(buff).Encode(ch)
	if err == nil {
		var resp *http.Response
		var req *http.Request
		req, err = http.NewRequest("POST", u.String(), buff)
		if err == nil {
			req.Header.Set("Content-Type", "application/json; encoding=UTF-8")
			resp, err = c.http.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == 200 {
					var chnlResp ChannelResponse
					err = json.NewDecoder(resp.Body).Decode(&chnlResp)
					if err == nil {
						if chnlResp.Meta.Code == 200 {
							chnl = &chnlResp.Data
						} else {
							err = ErrInvalidBackendResponse
						}
					}
				} else {
					err = ErrInvalidBackendResponse
				}
			}
		}
	}
	return
}
