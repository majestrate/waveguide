package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"waveguide/lib/adn"
	"waveguide/lib/config"
	"waveguide/lib/log"
	"waveguide/lib/model"
)

var sessionStore = sessions.NewCookieStore(config.GetAPISecret())

const sessionName = "waveguided-session"
const sessionKeyUserID = "user-id"
const sessionKeyToken = "user-token"
const sessionKeyUserName = "user-name"
const sessionKeyAvatarURL = "user-avatar-url"

func (r *Routes) GetCurrentUser(c *gin.Context) *model.UserInfo {

	var uid, token, name, avatar string
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		v := s.Values[sessionKeyUserID]
		if v != nil {
			uid = fmt.Sprintf("%s", v)
		}
		v = s.Values[sessionKeyToken]
		if v != nil {
			token = fmt.Sprintf("%s", v)
		}
		v = s.Values[sessionKeyUserName]
		if v != nil {
			name = fmt.Sprintf("%s", v)
		}
		v = s.Values[sessionKeyAvatarURL]
		if v != nil {
			avatar = fmt.Sprintf("%s", v)
		}
	}
	u := &model.UserInfo{
		UserID:    uid,
		Token:     token,
		Name:      name,
		AvatarURL: avatar,
	}
	return u
}

func (r *Routes) ResetUser(c *gin.Context) {
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		s.Values[sessionKeyUserID] = ""
		s.Values[sessionKeyToken] = ""
		s.Values[sessionKeyUserName] = ""
		s.Values[sessionKeyAvatarURL] = ""
		s.Save(c.Request, c.Writer)
	} else {
		log.Warnf("failed to reset user: %s", err.Error())
	}
}

func (r *Routes) SetCurrentUser(u adn.User, c *gin.Context) {
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		log.Infof("set user object: %q", u)
		s.Values[sessionKeyUserID] = u.ID
		s.Values[sessionKeyToken] = u.Token
		s.Values[sessionKeyUserName] = u.Username
		s.Values[sessionKeyAvatarURL] = u.Avatar.URL
		s.Save(c.Request, c.Writer)
	} else {
		log.Errorf("Failed to get session: %s", err.Error())
	}
}

func (r *Routes) CurrentUserLoggedIn(c *gin.Context) bool {
	return r.GetCurrentUser(c).UserID != ""
}

func (r *Routes) ServeUser(c *gin.Context) {
	// TODO: implement
}
