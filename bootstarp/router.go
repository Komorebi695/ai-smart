package bootstarp

import (
	"ai-smart/controller"
	"ai-smart/global"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func RegisterPath(s *gin.Engine) {
	s.Use(controller.CorsMiddleware)

	if global.Config.App.Env != "prod" {
		s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world!")
	})
	//controller

	gpt := s.Group("/v1/chat")
	{
		con := controller.NewOpenAIController()
		registerPostPath(gpt, "/completions", con.Completions)
	}

}

func registerPostPath(r *gin.RouterGroup, path string, f interface{}, middle ...gin.HandlerFunc) {
	if middle == nil {
		middle = make([]gin.HandlerFunc, 0)
	}

	middle = append(middle, GinHandelWrap2(f))
	r.POST(path, middle...)
}
