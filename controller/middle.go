package controller

import "github.com/gin-gonic/gin"

func ginJsonResponse(c *gin.Context, data interface{}) {
	c.JSON(200, data)
	c.Abort()
}
