package frontend

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"waveguide/lib/config"
	"waveguide/lib/model"
)

var sessionStore = sessions.NewCookieStore(config.GetAPISecret())

const sessionName = "waveguided-session"
const sessionKeyUserID = "user-id"

func (r *Routes) GetCurrentUser(c *gin.Context) *model.UserInfo {
	var uid int64
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		uid, _ = s.Values[sessionKeyUserID].(int64)
	}
	return &model.UserInfo{
		UserID: uid,
	}
}

func (r *Routes) SetCurrentUser(uid int64, c *gin.Context) {
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		s.Values[sessionKeyUserID] = uid
		s.Save(c.Request, c.Writer)
	}
}

func (r *Routes) CurrentUserLoggedIn(c *gin.Context) bool {
	// return r.GetCurrentUser(c).UserID != 0
	return true
}

func (r *Routes) ServeUser(c *gin.Context) {
	// TODO: implement
}
