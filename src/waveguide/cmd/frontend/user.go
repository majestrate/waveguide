package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"waveguide/lib/config"
	"waveguide/lib/log"
	"waveguide/lib/model"
	"waveguide/lib/oauth"
)

var sessionStore = sessions.NewCookieStore(config.GetAPISecret())

const sessionName = "waveguided-session"
const sessionKeyUserID = "user-id"
const sessionKeyToken = "user-token"
const sessionKeyUserName = "user-name"

func (r *Routes) GetCurrentUser(c *gin.Context) *model.UserInfo {
	var uid, token, name string
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
	}
	return &model.UserInfo{
		UserID: uid,
		Token:  token,
		Name:   name,
	}
}

func (r *Routes) SetCurrentUser(u oauth.User, c *gin.Context) {
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		log.Infof("set user object: %q", u)
		s.Values[sessionKeyUserID] = u.ID
		s.Values[sessionKeyToken] = u.Token
		s.Values[sessionKeyUserName] = u.Username
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
