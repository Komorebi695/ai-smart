package global

import (
	"ai-smart/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var App = new(Application)

type Application struct {
	Viper  *viper.Viper
	Config config.Configuration
	Log    *zap.Logger
}
