package model

import (
	"time"
	"waveguide/lib/model/util"
)

type UserInfo struct {
	UserID    string
	Token     string
	Name      string `form:"name" binding:"required"`
	Login     string
	Password  string `form:"password" binding:"required"`
	AvatarURL string
	Updated   time.Time
}

// UserCacheduration the amount of time between manditory profile info updates
const UserCacheDuration = time.Minute * 60

func (u UserInfo) Expired() bool {
	return time.Now().After(u.Updated.Add(UserCacheDuration))
}

func (u *UserInfo) Update() {
	u.Updated = time.Now()
}

func (u *UserInfo) CheckLogin(passwd string) bool {
	return util.CheckLogin(u.Login, passwd)
}

func (u *UserInfo) ChangePassword(newpasswd string) {
	u.Login = util.NewPassword(newpasswd)
}
