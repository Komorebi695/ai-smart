package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type MysqlInit struct {
	Master Mysql `json:"master"`
	Slaver Mysql `json:"slaver"`
}

type Mysql struct {
	Path         string `json:"path"`           // 服务器地址
	Port         string `json:"port"`           // 端口
	Config       string `json:"config"`         // 高级配置
	Dbname       string `json:"db-name"`        // 数据库名
	Username     string `json:"username"`       // 数据库用户名
	Password     string `json:"password"`       // 数据库密码
	TablePre     string `json:"table-pre"`      // 表前缀
	MaxIdleConns int    `json:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `json:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `json:"log-mode"`       // 是否开启Gorm全局日志
	LogZap       bool   `json:"log-zap"`        // 是否通过zap写入日志文件
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func InitDB(env, serviceName string, dbList []string) {

}

// initMysqlByConfig 初始化Mysql数据库用过传入配置
func initMysqlByConfig(m Mysql, gormConf gorm.Config) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gormConf); err != nil {
		log.Fatalf("gorm.Open err - config:%+v,gomConf:%+v", m, gormConf)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		sqlDB.SetConnMaxIdleTime(time.Minute)
		return db
	}
}

//func getMysqlConfig(env, serviceName, dbName, mode string, testMode bool, gormConf gorm.Config) (Mysql, gorm.Config) {
//	var confPath string
//	if testMode {
//		// 单元测试初始化模块，位于internal/test/base_test.go
//		confPath = fmt.Sprintf("../../conf/%s/mysql/%s.yaml", env, dbName)
//	} else {
//		// 项目初始化模块，位于根目录xxx.go
//		confPath = fmt.Sprintf(".conf/%s/mysql/%s.yaml", env, dbName)
//	}
//	var a MysqlInit
//
//}
