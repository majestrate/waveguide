package frontend

import (
	"github.com/gin-gonic/gin"
	"waveguide/lib/model"
)

func (r *Routes) ApiComment(c *gin.Context) {
	var comment model.Comment
	var has bool
	comment.Text, has = c.GetPostForm("comment")
	if has {

	}
}

func (r *Routes) ApiStreamComments(c *gin.Context) {

}
