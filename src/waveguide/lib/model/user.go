package model

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserID   int64
	Name     string
	Email    string `form:"email" binding:"required"`
	Login    string
	Password string `form:"password" binding:"required"`
}

// CheckLogin checks password against login credential
// returns true if the password matches otherwise false
func (u *UserInfo) CheckLogin() bool {
	hashed, _ := base64.StdEncoding.DecodeString(u.Login)
	return bcrypt.CompareHashAndPassword(hashed, []byte(u.Password)) == nil
}

func (u *UserInfo) UpdatePassword() {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Login = base64.StdEncoding.EncodeToString(hashed)
}
