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
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func (r *Routes) ApiAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if r.CurrentUserLoggedIn(c) {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func (r *Routes) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
