package frontend

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequiresCaptchaMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		solution, _ := c.GetPostForm("captcha_solution")
		id, _ := c.GetPostForm("captcha_id")
		if captcha.VerifyString(solution, id) {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"Error": "Bad Captcha",
			})
		}
	}
}

func (r *Routes) ApiAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if r.CurrentUserLoggedIn(c) {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"Error": "Not Authenticated",
			})
		}
	}
}
