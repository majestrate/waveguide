package model

import (
	"waveguide/lib/model/util"
)

type UserInfo struct {
	UserID   string
	Name     string `form:"name" binding:"required"`
	Login    string
	Password string `form:"password" binding:"required"`
}

func (u *UserInfo) CheckLogin(passwd string) bool {
	return util.CheckLogin(u.Login, passwd)
}

func (u *UserInfo) ChangePassword(newpasswd string) {
	u.Login = util.NewPassword(newpasswd)
}
