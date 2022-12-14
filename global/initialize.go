package global

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"main/internal/logging"
	"main/internal/settings"
	"main/middleware/dao"
	"os"
	"time"
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
	} else if env == "test" {
		Config = settings.InitConfig("conf/config_test.toml")
	} else {
		Config = settings.InitConfig("conf/config_local.toml")
	}

	fmt.Println("current deployParty:", DeployParty())
	if DeployParty() != DEPLOY_PARTY_P1 {
		fmt.Println("P1 Url:", P1Url())
	}
}

func DeployParty() string {
	deployParty := os.Getenv("DEPLOY_PARTY")
	//fmt.Println("current deployParty:", deployParty)
	if deployParty != DEPLOY_PARTY_P0 && deployParty != DEPLOY_PARTY_P1 && deployParty != DEPLOY_PARTY_BOTH {
		panic(errors.New("unsupported DEPLOY_PARTY"))
	}
	return deployParty
}

func P1Url() string {
	deployParty := os.Getenv("DEPLOY_PARTY")
	p1Url := os.Getenv("P1_URL")
	//fmt.Println("current deployParty:", deployParty)
	if deployParty != DEPLOY_PARTY_P1 && len(p1Url) == 0 {
		panic(errors.New("unsupported P1Url"))
	}
	return p1Url
}

func InitLogger() {

	if Config.Base.Env == "dev" {
		Logger = logging.InitLogger(logging.WithLogPath(Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
	} else {
		// k8s 直接输出到终端
		Logger = logging.InitLogger(logging.WithLogPath(Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
	}
}

type DbConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  int
}

func initDBEngine(config *DbConfig) *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName)
	println("connecting to DB, settings=", dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
	}), &gorm.Config{})

	if err != nil {
		panic("connect db error: " + err.Error())
	}
	sqlDb, _ := db.DB()
	if err := sqlDb.Ping(); err != nil {
		panic("DB ping error: " + err.Error())
	}

	sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(config.MaxLifetime))
	sqlDb.SetMaxIdleConns(config.MaxIdleConns)
	sqlDb.SetMaxOpenConns(config.MaxOpenConns)

	if Config.Base.Env != "prod" {
		db = db.Debug()
	}
	return db
}

func InitDB() {
	db := initDBEngine(&DbConfig{
		Host:         os.Getenv("MYSQL_HOST"),
		Port:         Config.DB.Port,
		User:         Config.DB.UserName,
		Password:     Config.DB.Password,
		DBName:       Config.DB.DBName,
		MaxIdleConns: Config.DB.MaxIdleConns,
		MaxOpenConns: Config.DB.MaxOpenConns,
		MaxLifetime:  Config.DB.MaxLifetime,
	})
	dao.SetDbEngine(db)
}

func InitAll() {
	InitConfig()
	InitLogger()
	InitDB()
}
