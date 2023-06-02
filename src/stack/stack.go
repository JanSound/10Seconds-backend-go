package stack

import "github.com/gin-gonic/gin"

func StackBeat(c *gin.Context) {
	c.JSON(200, gin.H{})
}
