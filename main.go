package main

import (
	"fmt"
	"main/global"
	"main/routes"
	"net/http"
	"os"
	"os/signal"
)

func closeResource() {
	// mysql 和 mongo 无需主动关闭资源

	// redis 连接池关闭
	//redis.CloseRedisConnection()
}

func init() {
	global.InitConfig()
	global.InitLogger()
}

func main() {
	//rnd, _ := new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	//b := rnd.Bytes() // big endian byte
	//reader := bytes.NewReader(b)
	//publicKey, privateKey, _ := eddsa.GenerateKey(reader)
	//fmt.Println(toLittleEdian(publicKey), privateKey)

	//clientKeypair, keyAgg := p0.KeyGen()
	//println("clientKeypair=", clientKeypair)
	//
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

	fmt.Println(fmt.Sprintf("Server Listen: http://0.0.0.0%v", global.Config.Base.Port))
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
