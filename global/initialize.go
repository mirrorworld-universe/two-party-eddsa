package global

import (
	"fmt"
	"main/internal/settings"
	"main/middleware/logging"
	"os"
)

func InitConfig() {
	// 根据环境变量读取不同的配置文件
	env := os.Getenv("ENV")
	fmt.Println("current env:", env)
	if env == "dev" {
		Config = settings.InitConfig("conf/config_dev.toml")
	} else if env == "staging" {
		Config = settings.InitConfig("conf/config_staging.toml")
	} else if env == "prod" {
		Config = settings.InitConfig("conf/config_prod.toml")
	} else if env == "uat" {
		Config = settings.InitConfig("conf/config_uat.toml")
	} else {
		Config = settings.InitConfig("conf/config_local.toml")
	}
}

func InitLogger() {

	if Config.Base.Env == "dev" {
		Logger = logging.InitLogger(logging.WithLogPath(Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
	} else {
		// k8s 直接输出到终端
		Logger = logging.InitLogger(logging.WithLogPath(Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
	}
}

//func initDB() {
//	sql.InitDBEngine(&sql.DbConfig{
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
