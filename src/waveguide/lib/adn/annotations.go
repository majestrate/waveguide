package adn

import (
	"encoding/json"
	"net/http"
	"net/url"
	"waveguide/lib/util"
)

const StreamAnnotation = "tv.gitgud"

type Stream struct {
	Online bool `json:"streaming"`
}

type Annotation struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (c *Client) GetAnnotations(token string, uid UID) (a []Annotation, err error) {
	var u *User
	u, err = c.GetUserByID(token, uid)
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
	buff := new(util.Buffer)
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
