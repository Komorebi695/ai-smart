package main

import (
	"ai-smart/bootstarp"
	"ai-smart/global"
	"ai-smart/initialize"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// 初始化配置
	initialize.InitConfig()
	global.App.Log = initialize.InitLog()
	global.App.Log.Info("zap log init success!")
	//config.InitDB("env", "local", []string{"main"}, false, gorm.Config{})

	s := gin.Default()
	// 注册路由
	bootstarp.RegisterPath(s)

	// 监听
	if err := endless.ListenAndServe(fmt.Sprintf(":%v", global.App.Config.App.Port), s); err != nil {
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
		log.Fatalf("unsupport method")
	}
}
