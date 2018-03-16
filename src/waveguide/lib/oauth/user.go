package oauth

import (
	"fmt"
	"waveguide/lib/model"
)

type User struct {
	ID          string       `json:"id"`
	Username    string       `json:"username"`
	Token       string       `json:"-"`
	Avatar      Image        `json:"avatar_image"`
	Cover       Image        `json:"cover_image"`
	Annotations []Annotation `json:"annotations"`
}

func (u User) ToModel() *model.UserInfo {
	return &model.UserInfo{
		UserID:    fmt.Sprintf("%d", u.ID),
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
