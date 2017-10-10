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
