package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"strconv"
	"time"
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
const sessionKeyAvatarURL = "user-avatar-url"
const sessionKeyLastUpdate = "user-updated"

func (r *Routes) GetCurrentUser(c *gin.Context) *model.UserInfo {

	var uid, token, name, avatar string
	var updated int64
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
		v = s.Values[sessionKeyLastUpdate]
		if v != nil {
			updated, err = strconv.ParseInt(fmt.Sprintf("%d", v), 10, 64)
			if err != nil {
				updated = 0
			}
		}
	}
	u := &model.UserInfo{
		UserID:    uid,
		Token:     token,
		Name:      name,
		AvatarURL: avatar,
		Updated:   time.Unix(updated, 0),
	}
	// refresh expired user info as needed
	if u.Expired() && r.oauth != nil && u.Token != "" {
		log.Infof("updating user info for %s", u.Name)
		user, err := r.oauth.GetUser(u.Token)
		if err == nil {
			r.SetCurrentUser(*user, c)
			u = user.ToModel()
			u.Update()
		} else {
			log.Warnf("failed to update user %s, %s", u.Name, err.Error())
		}
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
		s.Values[sessionKeyLastUpdate] = 0
		s.Save(c.Request, c.Writer)
	} else {
		log.Warnf("failed to reset user: %s", err.Error())
	}
}

func (r *Routes) SetCurrentUser(u oauth.User, c *gin.Context) {
	s, err := sessionStore.Get(c.Request, sessionName)
	if err == nil {
		log.Infof("set user object: %q", u)
		s.Values[sessionKeyUserID] = u.ID
		s.Values[sessionKeyToken] = u.Token
		s.Values[sessionKeyUserName] = u.Username
		s.Values[sessionKeyAvatarURL] = u.Avatar.URL
		s.Values[sessionKeyLastUpdate] = time.Now().UTC()
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
