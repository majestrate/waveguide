package model

type UserInfo struct {
	UserID   int64
	Name     string
	Email    string `form:"email" binding:"required"`
	Login    string
	Password string `form:"password" binding:"required"`
}
