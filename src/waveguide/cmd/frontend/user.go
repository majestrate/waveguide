package frontend

import (
	"github.com/gin-gonic/gin"
	"waveguide/lib/model"
)

func (r *Routes) GetCurrentUser(c *gin.Context) *model.UserInfo {
	// TODO: implement
	return &model.UserInfo{
		UserID: 0,
		Name:   "Anonymous",
	}
}

func (r *Routes) CurrentUserLoggedIn(c *gin.Context) bool {
	return r.GetCurrentUser(c).UserID != 0
}

func (r *Routes) ServeUser(c *gin.Context) {
	// TODO: implement
}
