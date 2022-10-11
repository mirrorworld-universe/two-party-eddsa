package controller

import (
	"github.com/gin-gonic/gin"
	"main/global"
	"main/internal/base_resp"
	"main/internal/binding"
	"main/internal/eddsa"
	"main/model/rest"
	"main/service/p0"
	"main/utils"
	"math/big"
	"net/http"
	"os"
)

func P0KeyGenRound1(c *gin.Context) {
	reqBody := rest.P0KeygenRound1Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}
	clientKeypair := new(eddsa.Keypair)
	keyAgg := new(eddsa.KeyAgg)

	var serverSKSeed *big.Int
	env := os.Getenv("ENV")
	if len(reqBody.ServerSKSeed) > 0 {
		if env == "prod" {
			global.Logger.Error("serverSKSeed should not be set in prod")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		serverSKSeed, _ = new(big.Int).SetString(reqBody.ServerSKSeed, 10)
	}

	if len(reqBody.ClientSKSeed) > 0 {
		clientSKSeed, ok := new(big.Int).SetString(reqBody.ClientSKSeed, 10)
		if !ok {
			global.Logger.Error("wrong clientSKSeed")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if serverSKSeed != nil {
			clientKeypair, keyAgg = p0.KeyGenRound1FromBothSeed(clientSKSeed, serverSKSeed)
		} else {
			clientKeypair, keyAgg = p0.KeyGenRound1FromSeed(clientSKSeed)
		}
	} else {
		clientKeypair, keyAgg = p0.KeyGenRound1NoSeed()
	}

	clientPubkeyByte := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPubkeyByte)
	println("keyagg=", keyAgg.ToString())
	resp := rest.P0KeygenRound1Response{
		ClientPubkeyBN: new(big.Int).SetBytes(clientPubkeyByte[:]).String(),
		KeyAgg:         keyAgg.Apk.ToHexString(),
	}
	base_resp.JsonResponseSimple(c, resp)
}

func P0SignRound1(c *gin.Context) {
	reqBody := rest.P0SignRound1Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}

	// this should be read from db by userId.
	clientSKSeed, _ := new(big.Int).SetString(reqBody.ClientSKSeed, 10)
	clientKeypair, keyAgg := p0.KeyGenRound1FromSeed(clientSKSeed)

	msgHash := utils.StringToBigInt(&reqBody.Msg)
	msgHashBN := msgHash.String()
	p0.SignRound1(&msgHashBN, clientKeypair, keyAgg)
	println("a")
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
