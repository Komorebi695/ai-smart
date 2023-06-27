package config

import (
	"ai-smart/initialize"
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
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
	return fmt.Sprintf("%v:%v@tcp(%v:%d)/%v?%v", m.Username, m.Password, m.Path, m.Port, m.Dbname, m.Config)
}

func IsGormFound(err error) error {
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func InitDB(env, serviceName string, dbList []string, testMode bool, gomConf gorm.Config) {
	passEnv := map[string]bool{
		"dev":  true,
		"prod": true,
		"test": true,
	}
	if _, ok := passEnv[env]; ok {
		log.Fatalf("InitDB env fail - serviceName:%s,dbList:%+V,gormConf:%+v", serviceName, dbList, gomConf)
	}
	DbMapsInit(env, serviceName, dbList, testMode, gomConf)
}

func DbMapsInit(env, serviceName string, initLists []string, testMode bool, gomConf gorm.Config) map[string]map[string]*gorm.DB {
	DbMaps := make(map[string]map[string]*gorm.DB)
	for _, v := range initLists {
		DbMaps[v] = map[string]*gorm.DB{
			"master": gormInit(env, serviceName, v, "master", testMode, gomConf),
			"slaver": gormInit(env, serviceName, v, "slaver", testMode, gomConf),
		}
	}
	return DbMaps
}

func gormInit(env, serviceName, dbName, mode string, test bool, gormConf gorm.Config) *gorm.DB {
	return initMysqlByConfig(getMysqlConfig(env, serviceName, dbName, mode, test, gormConf))
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

func getMysqlConfig(env, serviceName, dbName, mode string, testMode bool, gormConf gorm.Config) (Mysql, gorm.Config) {
	var confPath string
	if testMode {
		// 单元测试初始化模块，位于internal/test/base_test.go
		confPath = fmt.Sprintf("../../conf/%s/mysql/%s.yaml", env, dbName)
	} else {
		// 项目初始化模块，位于根目录xxx.go
		confPath = fmt.Sprintf(".conf/%s/mysql/%s.yaml", env, dbName)
	}
	confPath = "./config.yaml"
	var a MysqlInit
	file, err := os.Open(confPath)
	if err != nil {
		return Mysql{}, gorm.Config{}
	}
	if err := yaml.NewDecoder(file).Decode(&a); err != nil {
		return Mysql{}, gorm.Config{}
	}

	gormConf.Logger = DBLog{
		env:         env,
		serviceName: serviceName,
	}

	LogModeMap := map[string]logger.LogLevel{
		"dev":  3,
		"test": 2,
		"prod": 1,
	}
	gormConf.Logger = gormConf.Logger.LogMode(LogModeMap[env])
	if mode == "w" || mode == "write" || mode == "master" {
		return a.Master, gormConf
	} else {
		return a.Slaver, gormConf
	}
}

type DBLog struct {
	env         string
	serviceName string
	logLevel    logger.LogLevel
}

func (dBLog DBLog) Info(c context.Context, s string, i ...interface{}) {
	if dBLog.logLevel > 2 {
		log.Printf("sql:s:%+v,i:%+v", s, i)
	}
}

func (dBLog DBLog) Warn(c context.Context, s string, i ...interface{}) {
	if dBLog.logLevel > 2 {
		log.Printf("sql:Warn:s:%+v,i:%+v", s, i)
	}
}

func (dBLog DBLog) Error(c context.Context, s string, i ...interface{}) {
	if dBLog.logLevel > 2 {
		log.Printf("sql:Error:s:%+v,i:%+v", s, i)
	}
}

func (dBLog DBLog) Trace(c context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if err == nil {
		if dBLog.logLevel > 1 {
			sql, rowsAffected := fc()
			log.Printf("sql:%s", fmt.Sprintf("%+v,rowsAffected:%+v", sql, rowsAffected))
		}
	} else {
		log.Printf("sql:err:%s", err.Error())
	}
}

func (dBLog DBLog) LogMode(logLevel logger.LogLevel) logger.Interface {
	return DBLog{
		env:         dBLog.env,
		serviceName: dBLog.serviceName,
		logLevel:    logLevel,
	}
}

// DB 拿出DB
func DB(dbName string, mode ...string) *gorm.DB {
	if len(mode) > 0 {
		dbChange := strings.ToLower(mode[0])
		if dbChange == "w" || dbChange == "write" || dbChange == "master" {
			return initialize.DbMaps[dbName]["master"]
		} else {
			return initialize.DbMaps[dbName]["slaver"]
		}
	} else {
		return initialize.DbMaps[dbName]["slaver"]
	}
}
