package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"waveguide/lib/config"
)

type Client struct {
	conf config.OAuthConfig
}

func (c *Client) AuthURL(callback string) string {
	u, _ := url.Parse(c.conf.Provider + "oauth/authenticate?client_id=" + c.conf.ClientID + "&response_type=code")
	q := u.Query()
	q.Set("redirect_uri", callback)
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) GetUser(code, callback string) (user *User, err error) {
	var postdata url.Values
	postdata.Set("client_id", c.conf.ClientID)
	postdata.Set("client_secret", c.conf.ClientSecret)
	postdata.Set("grant_type", "authorization_code")
	postdata.Set("redirect_uri", callback)
	postdata.Set("code", code)
	var buff bytes.Buffer
	io.WriteString(&buff, postdata.Encode())
	var resp *http.Response
	resp, err = http.Post(c.conf.Provider+"oauth/acesss_token", "application/x-www-form-urlencoded", &buff)
	if err == nil {
		var u User
		err = json.NewDecoder(resp.Body).Decode(&u)
		if err == nil {
			user = &u
		}
	}
	return
}

func NewClient(c config.OAuthConfig) *Client {
	return &Client{
		conf: c,
	}
}
