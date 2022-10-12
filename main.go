package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"main/global"
	"main/internal/logging"
	"main/internal/settings"
	"main/middleware/dao"
	validator2 "main/middleware/validator"
	"main/model/db"
	"main/routes"
	"main/utils"
	"math/big"
	"net/http"
	"os"
	"os/signal"
)

func closeResource() {
	// mysql 和 mongo 无需主动关闭资源

	// redis 连接池关闭
	//redis.CloseRedisConnection()
}

func InitConfig() {
	// 根据环境变量读取不同的配置文件
	env := os.Getenv("ENV")
	fmt.Println("current env:", env)
	if env == "dev" {
		global.Config = settings.InitConfig("conf/config_dev.toml")
	} else if env == "staging" {
		global.Config = settings.InitConfig("conf/config_staging.toml")
	} else if env == "prod" {
		global.Config = settings.InitConfig("conf/config_prod.toml")
	} else {
		global.Config = settings.InitConfig("conf/config_local.toml")
	}
}

func InitLogger() {

	if global.Config.Base.Env == "dev" {
		global.Logger = logging.InitLogger(logging.WithLogPath(global.Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
	} else {
		// k8s 直接输出到终端
		global.Logger = logging.InitLogger(logging.WithLogPath(global.Config.Log.Path), logging.WithOutput(logging.ONLY_TERMINAL))
	}
}

func InitDB() {
	dao.InitDBEngine(&dao.DbConfig{
		Host:         global.Config.DB.Host,
		Port:         global.Config.DB.Port,
		User:         global.Config.DB.UserName,
		Password:     global.Config.DB.Password,
		DBName:       global.Config.DB.DBName,
		MaxIdleConns: global.Config.DB.MaxIdleConns,
		MaxOpenConns: global.Config.DB.MaxOpenConns,
		MaxLifetime:  global.Config.DB.MaxLifetime,
	})

	// auto migrate
	dao.GetDbEngine().AutoMigrate(&db.MPCWallet{})
}

func test(b string) {
	byte, _ := new(big.Int).SetString(b, 10)
	println(utils.BytesToStr(byte.Bytes()))
}

func main() {
	// ServerSignSecondMsgR=
	//test("64483640907921859330156930639679621060811220324382896999100840459413651822902")
	//ServerSigRBN=
	//test("13501676355858580285518499925993993426720953276435438562170709642663523246311")
	//ServerSigSmallSBN=
	//test("56362407394169897738728800006116437916583347118160729121811866192132326280458")
	//rnd, _ := new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	//b := rnd.Bytes() // big endian byte
	//reader := bytes.NewReader(b)
	//publicKey, privateKey, _ := eddsa.GenerateKey(reader)
	//fmt.Println(toLittleEdian(publicKey), privateKey)

	//clientKeypair, keyAgg := p0.KeyGen()
	//println("clientKeypair=", clientKeypair)
	////
	//println("\n\n************ SIGN now *************")
	//msg := "hello"
	//p0.Sign(&msg, clientKeypair, keyAgg)

	//R := "1dd9ad91d660e104fc02043e7bbe0c303f1bfc1c012689ab8c2d38c4ae6be0e7"
	//s := "7bf0d2eb8027a65988c43a4c79e70f3ab67eadf1a8a852b5cf34ef1ace192407"
	//pubkey := "790c23f4a2f065fa4cebf77a005f75ad7a528c8de4ca64e4e5c681c17663514e"
	//p0.Verify(&msg, &R, &s, &pubkey)
	//bn, _ := new(big.Int).SetString("2D282E87852DE00E981EFAAD08A9F435CF41CDFD7BD9EB3DBFBE08AE5537160", 16)
	//serverKeypair := eddsa.CreateKeyPairFromSeed(bn)
	//println("serverKeypair=", serverKeypair.ToString())
	//
	//p1.Sign(serverKeypair, keyAgg)

	InitConfig()
	InitLogger()
	InitDB()

	// custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validbn", validator2.ValidBN)
	}

	router := routes.NewRouter()
	srv := &http.Server{
		Addr:    global.Config.Base.Port,
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			global.Logger.Error("Gin server start error:", err.Error())
			panic(err.Error())
		}
	}()

	fmt.Println(fmt.Sprintf("Server Listen: rest://0.0.0.0%v", global.Config.Base.Port))
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
