package config

type Configuration struct {
	App App `mapstructure:"app" json:"app" yaml:"app"`
	Log Log `mapstructure:"log" json:"log" yaml:"log"`
	Jwt Jwt `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}

type App struct {
	Env     string `json:"env" yaml:"env" mapstructure:"env"`
	Port    string `json:"port" yaml:"port" mapstructure:"port"`
	AppName string `json:"app_name" yaml:"app_name" mapstructure:"app_name"`
	AppUrl  string `json:"app_url" yaml:"app_url" mapstructure:"app_url"`
}
