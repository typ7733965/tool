package db

import (
	"fmt"
	"github.com/typ7733965/tool/comfunc/env"
	"github.com/typ7733965/tool/config"
	"github.com/typ7733965/tool/enum"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"log"
	"net/url"
	"os"
	"time"
)

func NewMysql(config *config.MysqlConfig) (dbInst *gorm.DB, err error) {
	dsn := config.Master.Dsn
	if dsn == "" {
		dsn = parseDsn(&config.Master)
	}
	level := glogger.Info
	slow := time.Second * 5
	if env.GetAppInfo().Stage != enum.EvnStageLocal {
		level = glogger.Error
	}
	if config.SlowThreshold > 0 {
		slow = time.Second * time.Duration(config.SlowThreshold)
	}
	if config.LogLevel > 0 {
		level = glogger.LogLevel(config.LogLevel)
	} else {
		stage := os.Getenv(enum.StageKey)
		if stage == enum.EvnStageProduction || stage == "master" {
			// 生产环境warn级别
			level = glogger.Warn
		}
	}
	newLogger := glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		glogger.Config{
			SlowThreshold:             slow,  // Slow SQL threshold
			LogLevel:                  level, // Log level
			IgnoreRecordNotFoundError: true,  // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,  // Don't include params in the SQL log
			Colorful:                  false, // Disable color
		},
	)
	dbInst, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: !config.PluralTable}, Logger: newLogger})
	if err != nil {
		return
	}
	var replicas []gorm.Dialector
	if len(config.Slave) > 0 {
		for i, _ := range config.Slave {
			slaveDsn := parseDsn(&config.Slave[i])
			replicas = append(replicas, mysql.Open(slaveDsn))
		}
	}
	err = dbInst.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(config.MaxIdleTime).
		SetConnMaxLifetime(config.MaxLifeTime).
		SetMaxIdleConns(config.MaxIdleConn).
		SetMaxOpenConns(config.MaxOpenConn))
	if err != nil {
		return
	}
	//测试连接
	sqlDB, err := dbInst.DB()
	if err != nil {
		return
	}
	err = sqlDB.Ping()
	return
}
func parseDsn(conf *config.MysqlSingleConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s&timeout=30s",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.Db,
		url.QueryEscape("UTC+8"))
	return dsn
}
