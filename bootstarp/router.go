package bootstarp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterPath(s *gin.Engine) {
	s.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world!")
	})
	//controller

}

func registerPostPath(r *gin.RouterGroup, path string, f interface{}, middle ...gin.HandlerFunc) {
	if middle == nil {
		middle = make([]gin.HandlerFunc, 0)
	}

	middle = append(middle, GinHandelWrap2(f))
	r.POST(path, middle...)
}
