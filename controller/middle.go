package controller

import "github.com/gin-gonic/gin"

func ginJsonResponse(c *gin.Context, data interface{}) {
	c.JSON(200, data)
	c.Abort()
}

// CorsMiddleware 跨域
func CorsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Cookie, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,authorization, Ad-Code, ad-code, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
