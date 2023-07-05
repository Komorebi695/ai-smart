package global

import (
	"ai-smart/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//var App = new(Application)

var (
	DbMaps map[string]map[string]*gorm.DB
)

//type Application struct {
//	Viper  *viper.Viper
//	Config config.Configuration
//	Log    *zap.Logger
//}

var (
	Viper  *viper.Viper
	Config config.Configuration
	Log    *zap.Logger
)
