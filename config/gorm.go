package config

import (
	"fmt"
)

type MySQl struct {
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	Dbname              string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Username            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Config              string `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

func (m *MySQl) Dsn() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%d)/%v?%v", m.Username, m.Password, m.Host, m.Port, m.Dbname, m.Config)
}

//func IsGormFound(err error) error {
//	if err == gorm.ErrRecordNotFound {
//		return nil
//	}
//	return err
//}
