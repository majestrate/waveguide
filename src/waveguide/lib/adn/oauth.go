package adn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"waveguide/lib/config"
	"waveguide/lib/httppool"
	"waveguide/lib/util"
)

var ErrInvalidBackendResponse = errors.New("invalid response from oauth backend")
var ErrFailedToContactBackend = errors.New("failed to contact oauth backend server")
var ErrOAuthClosing = errors.New("oauth backend is already closed")

type Client struct {
	conf    config.ADNConfig
	http    *httppool.Client
	closing bool
}

// implements io.Closer
func (c *Client) Close() error {
	if c.closing {
		return ErrOAuthClosing
	}
	c.closing = true
	return c.http.Close()
}

func (c *Client) doHTTP(method, path, token string, data interface{}) (resp *http.Response, err error) {
	if path[0] == '/' {
		path = path[1:]
	}
	buff := new(util.Buffer)
	if data != nil {
		err = json.NewEncoder(buff).Encode(data)
	}
	if err == nil {
		var req *http.Request
		if data == nil {
			req, err = http.NewRequest(method, c.conf.Provider+path, nil)
		} else {
			req, err = http.NewRequest(method, c.conf.Provider+path, buff)
		}
		if err == nil {
			if buff != nil {
				req.Header.Set("Content-Size", fmt.Sprintf("%d", buff.Len()))
				req.Header.Set("Content-Type", "application/json; encoding=UTF-8")
			}
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err = c.http.Do(req)
		}
	}
	return
}

func (c *Client) AuthURL(callback string) string {
	u, _ := url.Parse(c.conf.Provider + "oauth/authenticate?client_id=" + c.conf.ClientID + "&response_type=code")
	q := u.Query()
	q.Set("redirect_uri", callback)
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) GrantUser(code, callback string) (user *User, err error) {
	postdata := make(url.Values)
	postdata.Set("client_id", c.conf.ClientID)
	postdata.Set("client_secret", c.conf.ClientSecret)
	postdata.Set("grant_type", "authorization_code")
	postdata.Set("redirect_uri", callback)
	postdata.Set("code", code)
	buff := new(util.Buffer)
	io.WriteString(buff, postdata.Encode())
	var resp *http.Response
	var req *http.Request
	req, err = http.NewRequest("POST", c.conf.Provider+"oauth/access_token", buff)
	if err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err = c.http.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var tok TokenRequest
			err = json.NewDecoder(resp.Body).Decode(&tok)
			if err == nil {
				user = &User{
					ID:       tok.Token.User.ID,
					Username: tok.Token.User.Username,
					Token:    tok.AccessToken,
				}
			}
			resp.Body.Close()
		} else {
			err = ErrFailedToContactBackend
		}
	}
	return
}

func NewClient(c config.ADNConfig) *Client {
	return &Client{
		conf: c,
		http: httppool.New(c.Workers),
	}
}
