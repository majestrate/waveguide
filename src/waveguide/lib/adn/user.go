package adn

import (
	"encoding/json"
	"fmt"
	"waveguide/lib/model"
)

type User struct {
	ID          UID          `json:"id"`
	Username    string       `json:"username"`
	Token       string       `json:"-"`
	Avatar      Image        `json:"avatar_image"`
	Cover       Image        `json:"cover_image"`
	Annotations []Annotation `json:"annotations"`
}

type UserRequest struct {
	Data User         `json:"data"`
	Meta MetaResponse `json:"meta"`
}

func (u User) ToModel() *model.UserInfo {
	return &model.UserInfo{
		UserID:    u.ID.String(),
		Name:      u.Username,
		Token:     u.Token,
		AvatarURL: u.Avatar.URL,
	}
}

type Image struct {
	Height  int    `json:"height"`
	Width   int    `json:"width"`
	URL     string `json:"url"`
	Default bool   `json:"is_default"`
}

type Token struct {
	User User `json:"user"`
}

type TokenRequest struct {
	AccessToken string `json:"access_token"`
	Token       Token  `json:"token"`
}

func (c *Client) GetUserByID(token string, id UID) (user *User, err error) {
	resp, err := c.doHTTP("GET", fmt.Sprintf("/stream/0/users/%s?include_annotations=1", id.String()), token, nil)
	if err == nil {
		var userReq UserRequest
		err = json.NewDecoder(resp.Body).Decode(&userReq)
		resp.Body.Close()
		if err == nil {
			if userReq.Meta.Code == 200 {
				user = &userReq.Data
			} else {
				err = ErrInvalidBackendResponse
			}

		}
	}
	return
}

func (c *Client) GetCurrentUser(token string) (user *User, err error) {
	resp, err := c.doHTTP("GET", "/stream/0/token", token, nil)
	if err == nil {
		var tokenReq TokenInfoRequest
		err = json.NewDecoder(resp.Body).Decode(&tokenReq)
		resp.Body.Close()
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
	return
}
