package main

//
//import (
//	"fmt"
//	"main/global"
//	"main/internal/logging"
//	"main/internal/settings"
//	"main/middleware/dao"
//	"os"
//)
//
//func InitConfig() {
//	// 根据环境变量读取不同的配置文件
//	env := os.Getenv("ENV")
//	fmt.Println("current env:", env)
//	if env == "dev" {
//		global.Config = settings.InitConfig("conf/config_dev.toml")
//	} else if env == "staging" {
//		global.Config = settings.InitConfig("conf/config_staging.toml")
//	} else if env == "prod" {
//		global.Config = settings.InitConfig("conf/config_prod.toml")
//	} else {
//		global.Config = settings.InitConfig("conf/config_local.toml")
//	}
//}
//
//func InitLogger() {
//
//	if global.Config.Base.Env == "dev" {
//		global.Logger = logging.InitLogger(logging.WithLogPath(global.Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
//	} else {
//		// k8s 直接输出到终端
//		global.Logger = logging.InitLogger(logging.WithLogPath(global.Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
//	}
//}
//
//func InitDB() {
//	dao.InitDBEngine(&dao.DbConfig{
//		Host:         global.Config.DB.Host,
//		Port:         global.Config.DB.Port,
//		User:         global.Config.DB.UserName,
//		Password:     global.Config.DB.Password,
//		DBName:       global.Config.DB.DBName,
//		MaxIdleConns: global.Config.DB.MaxIdleConns,
//		MaxOpenConns: global.Config.DB.MaxOpenConns,
//		MaxLifetime:  global.Config.DB.MaxLifetime,
//	})
//}
