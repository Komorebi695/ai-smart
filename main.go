package main

import (
	"ai-smart/internal/controller"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	s := gin.Default()

	// 注册路由
	controller.RegisterPath(s)

	// 监听
	if err := endless.ListenAndServe(fmt.Sprintf(":%v", "8080"), s); err != nil {
		log.Fatalf("listen err:%v", err)
	}
}

func RegisterPath(r *gin.RouterGroup, method, path string, h http.HandlerFunc) {
	switch method {
	case http.MethodPost:
		r.POST(path, gin.WrapF(h))
	case http.MethodGet:
		r.GET(path, gin.WrapF(h))
	default:
		panic("unsupport method")
	}
}
