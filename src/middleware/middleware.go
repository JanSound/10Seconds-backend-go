package middleware

import "github.com/gin-gonic/gin"

func MyMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Do something before request
		c.Next()
		// Do something after request
	}
}
