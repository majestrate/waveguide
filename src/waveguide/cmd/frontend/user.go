package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"waveguide/lib/config"
	"waveguide/lib/model"
)

var sessionStore = sessions.NewCookieStore(config.GetAPISecret())

const sessionName = "waveguided-session"
const sessionKeyUserID = "user-id"

func (r *Routes) GetCurrentUser(c *gin.Context) *model.UserInfo {
	var uid string
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		uid = fmt.Sprintf("%s", s.Values[sessionKeyUserID])
	}
	return &model.UserInfo{
		UserID: uid,
	}
}

func (r *Routes) SetCurrentUser(uid string, c *gin.Context) {
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		s.Values[sessionKeyUserID] = uid
		s.Save(c.Request, c.Writer)
	}
}

func (r *Routes) CurrentUserLoggedIn(c *gin.Context) bool {
	return r.GetCurrentUser(c).UserID != ""
}

func (r *Routes) ServeUser(c *gin.Context) {
	// TODO: implement
}
