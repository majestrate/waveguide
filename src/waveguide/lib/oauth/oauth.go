package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"waveguide/lib/config"
	"waveguide/lib/model"
)

var ErrInvalidBackendResponse = errors.New("invalid response from oauth backend")
var ErrFailedToContactBackend = errors.New("failed to contact oauth backend server")

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

func (c *Client) SubmitComment(comment model.Comment) (err error) {
	u, _ := url.Parse(c.conf.Provider + fmt.Sprintf("channels/%s/messages", comment.ChannelID))
	q := u.Query()
	q.Set("access_token", comment.User.Token)
	u.RawQuery = q.Encode()
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(map[string]interface{}{
		"text": comment.Text,
	})
	if err == nil {
		var resp *http.Response
		resp, err = http.Post(u.String(), "application/json; encoding=UTF-8", buff)
		if err == nil {
			if resp.StatusCode == 200 {
				// TODO: check json response
			} else {
				err = ErrInvalidBackendResponse
			}
		}
	}
	return
}

func (c *Client) AnnounceStream(token, message string) (err error) {
	u, _ := url.Parse(c.conf.Provider + "channels/5/messages")
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(map[string]interface{}{
		"text": message,
	})
	if err == nil {
		var resp *http.Response
		resp, err = http.Post(u.String(), "application/json; encoding=UTF-8", buff)
		if err == nil {
			if resp.StatusCode == 200 {
				// TODO: check json response
			} else {
				err = ErrInvalidBackendResponse
			}
		}
	}
	return
}

func (c *Client) GetUser(token string) (user *User, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", c.conf.Provider+"stream/0/token", nil)
	if err == nil {
		req.Header.Set("Authorization", "Bearer "+token)
		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var tokenReq TokenInfoRequest
			err = json.NewDecoder(resp.Body).Decode(&tokenReq)
			if err == nil {
				user = &User{
					ID:       tokenReq.Data.User.ID,
					Username: tokenReq.Data.User.Username,
					Token:    token,
				}
			}
		}
	}
	return
}

func (c *Client) GrantUser(code, callback string) (user *User, err error) {
	postdata := make(url.Values)
	postdata.Set("client_id", c.conf.ClientID)
	postdata.Set("client_secret", c.conf.ClientSecret)
	postdata.Set("grant_type", "authorization_code")
	postdata.Set("redirect_uri", callback)
	postdata.Set("code", code)
	buff := new(bytes.Buffer)
	io.WriteString(buff, postdata.Encode())
	var resp *http.Response
	resp, err = http.Post(c.conf.Provider+"oauth/access_token", "application/x-www-form-urlencoded", buff)
	if err == nil {
		defer resp.Body.Close()
		var tok TokenRequest
		err = json.NewDecoder(resp.Body).Decode(&tok)
		if err == nil {
			user = &User{
				Token:    tok.AccessToken,
				Username: tok.Token.User.Username,
				ID:       tok.Token.User.ID,
			}
		} else {
			err = ErrInvalidBackendResponse
		}
		resp.Body.Close()
	} else {
		err = ErrFailedToContactBackend
	}
	return
}

func NewClient(c config.OAuthConfig) *Client {
	return &Client{
		conf: c,
	}
}
