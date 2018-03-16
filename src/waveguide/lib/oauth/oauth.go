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

func (c *Client) EnsureStreamChat(token string) (chatid int64, err error) {
	var chnl *Channel
	chnl, err = c.CreateChannel(token, StreamAnnotation)
	if err == nil {
		if chnl == nil {
			err = ErrInvalidBackendResponse
		} else {
			chatid = chnl.ID
		}
	}
	return
}

func (c *Client) DeleteChannel(token string, chnlID int64) (err error) {
	u, _ := url.Parse(c.conf.Provider + fmt.Sprintf("stream/0/channels/%d", chnlID))
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

func (c *Client) CreateChannel(token string, chanType string) (chnl *Channel, err error) {
	u, _ := url.Parse(c.conf.Provider + "stream/0/channels")
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(map[string]interface{}{
		"type": chanType,
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
				if resp.StatusCode == 200 {
					var chnlResp ChannelResponse
					err = json.NewDecoder(resp.Body).Decode(&chnlResp)
					if err == nil {
						if chnlResp.Meta.Code == 200 {
							chnl = &chnlResp.Data
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

func (c *Client) GetAnnotations(token string) (a []Annotation, err error) {
	var u *User
	u, err = c.GetUser(token)
	if err == nil && u != nil {
		a = u.Annotations
	}
	return
}

func (c *Client) putAnnotations(token string, a []Annotation) (err error) {
	u, _ := url.Parse(c.conf.Provider + "stream/0/users/me?include_annotations=1")
	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(map[string]interface{}{
		"annotations": a,
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
	var anos []Annotation
	anos, err = c.GetAnnotations(token)
	if err == nil {
		for _, a := range anos {
			if a.Type == StreamAnnotation {
				a.Value = Stream{
					Online: status,
				}
			}
		}
		err = c.putAnnotations(token, anos)
	}
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

func (c *Client) SubmitPost(token string, channel int64, post Post) (err error) {
	u, _ := url.Parse(c.conf.Provider + fmt.Sprintf("channels/%d/messages?include_annotations=1", channel))
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
	req, err = http.NewRequest("GET", c.conf.Provider+"stream/0/token?include_annotations=1", nil)
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
					ID:          tokenReq.Data.User.ID,
					Username:    tokenReq.Data.User.Username,
					Token:       token,
					Avatar:      tokenReq.Data.User.Avatar,
					Cover:       tokenReq.Data.User.Cover,
					Annotations: tokenReq.Data.User.Annotations,
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
