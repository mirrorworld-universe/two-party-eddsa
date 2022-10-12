package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/finder"
	"main/global"
	"main/internal/base_resp"
	"main/internal/binding"
	"main/internal/eddsa"
	error_code "main/internal/err_code"
	"main/model/db"
	"main/model/rest"
	"main/service/p0"
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

	var wallet *db.MPCWallet
	var err *error
	if len(reqBody.ClientSKSeed) > 0 {
		clientSKSeed, ok := new(big.Int).SetString(reqBody.ClientSKSeed, 10)
		if !ok {
			global.Logger.Error("wrong clientSKSeed")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if serverSKSeed != nil {
			clientKeypair, keyAgg, wallet, err = p0.KeyGenRound1FromBothSeed(clientSKSeed, serverSKSeed)
		} else {
			clientKeypair, keyAgg, wallet, err = p0.KeyGenRound1FromSeed(clientSKSeed)
		}
	} else {
		clientKeypair, keyAgg, wallet, err = p0.KeyGenRound1NoSeed()
	}

	bsp := error_code.NewBaseResp()
	if err != nil {
		bsp.SetMsg(error_code.InternalError, fmt.Sprintf("%v", err))
		base_resp.JsonResponse(c, bsp, nil)
		return
	}

	clientPubkeyByte := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPubkeyByte)

	println("keyagg=", keyAgg.ToString())
	resp := rest.P0KeygenRound1Response{
		ClientPubkeyBN: new(big.Int).SetBytes(clientPubkeyByte[:]).String(),
		KeyAgg:         keyAgg.Apk.ToHexString(),
		UserId:         wallet.UserId,
	}
	base_resp.JsonResponseSimple(c, resp)
}

func P0SignRound1(c *gin.Context) {
	reqBody := rest.P0SignRound1Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}

	wallet := finder.FindP0ByUserId(&reqBody.UserId)
	seedBN, _ := new(big.Int).SetString(wallet.SeedBN, 10)
	clientKeypair := eddsa.CreateKeyPairFromSeed(seedBN)

	apkBN, _ := new(big.Int).SetString(wallet.KeyAggAPKBN, 10)
	hashBN, _ := new(big.Int).SetString(wallet.KeyAggHashBN, 10)
	keyAgg := eddsa.NewKeyAggFromBNs(apkBN, hashBN)

	R, s, err := p0.SignRound1(&wallet.UserId, &reqBody.Msg, clientKeypair, keyAgg)
	bsp := error_code.NewBaseResp()
	if err != nil {
		bsp.SetMsg(error_code.InternalError, fmt.Sprintf("%v", err))
		base_resp.JsonResponse(c, bsp, nil)
		return
	}
	data := rest.P0SignRound1Response{
		R:      *R,
		SmallS: *s,
	}
	base_resp.JsonResponse(c, bsp, data)
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
