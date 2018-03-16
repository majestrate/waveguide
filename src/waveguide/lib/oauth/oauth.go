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
	"waveguide/lib/httppool"
	"waveguide/lib/model"
)

var ErrInvalidBackendResponse = errors.New("invalid response from oauth backend")
var ErrFailedToContactBackend = errors.New("failed to contact oauth backend server")
var ErrOAuthClosing = errors.New("oauth backend is already closed")

type Client struct {
	conf    config.OAuthConfig
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
		var req *http.Request
		req, err = http.NewRequest("POST", u.String(), buff)
		if err == nil {
			req.Header.Set("Content-Type", "application/json; encoding=UTF-8")
			resp, err = c.http.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					err = ErrInvalidBackendResponse
				}
			}
		}
	}
	return
}

func (c *Client) EnsureChat(token string) (chatid string, err error) {
	// TODO: implement this
	chatid = "5"
	return
}

func (c *Client) UpdateAnnotation(token string, a Annotation) (err error) {
	u, _ := url.Parse(c.conf.Provider + "stream/0/users/me")
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(map[string]interface{}{
		"annotations": []Annotation{a},
	})
	if err == nil {
		var resp *http.Response
		var req *http.Request
		req, err = http.NewRequest("PATCH", u.String(), buff)
		if err == nil {
			req.Header.Set("Content-Type", "application/json; encoding=UTF-8")
			resp, err = c.http.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					err = ErrInvalidBackendResponse
				}
			}
		}
	}
	return
}

func (c *Client) streamStatus(token, uid string, status bool) (err error) {
	err = c.UpdateAnnotation(token, Annotation{
		Type: StreamAnnotation,
		Value: Stream{
			Online: status,
		},
	})
	return
}

func (c *Client) StreamOnline(token, uid string) (err error) {
	err = c.streamStatus(token, uid, true)
	return
}

func (c *Client) StreamOffline(token, uid string) (err error) {
	err = c.streamStatus(token, uid, false)
	return
}

func (c *Client) SubmitPost(token, channel string, post Post) (err error) {
	u, _ := url.Parse(c.conf.Provider + fmt.Sprintf("channels/%s/messages?include_post_annotations=1", channel))
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(post)
	if err == nil {
		var resp *http.Response
		var req *http.Request
		req, err = http.NewRequest("POST", u.String(), buff)
		if err == nil {
			req.Header.Set("Content-Type", "application/json; encoding=UTF-8")
			resp, err = c.http.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					err = ErrInvalidBackendResponse
				}
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
		resp, err = c.http.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var tokenReq TokenInfoRequest
			err = json.NewDecoder(resp.Body).Decode(&tokenReq)
			if err == nil {
				user = &User{
					ID:       tokenReq.Data.User.ID,
					Username: tokenReq.Data.User.Username,
					Token:    token,
					Avatar:   tokenReq.Data.User.Avatar,
					Cover:    tokenReq.Data.User.Cover,
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
	}
	return
}

func NewClient(c config.OAuthConfig) *Client {
	return &Client{
		conf: c,
		http: httppool.New(c.Workers),
	}
}
