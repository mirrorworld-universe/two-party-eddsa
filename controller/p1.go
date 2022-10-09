package controller

import (
	"github.com/gin-gonic/gin"
	"main/global"
	"main/internal/base_resp"
	"main/internal/binding"
	"main/model/rest"
	"main/service/p1"
	"math/big"
	"net/http"
)

func P1KeyGenRound1(c *gin.Context) {
	reqBody := rest.P1KeygenRound1Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}
	clientPubkeyBN, ok := new(big.Int).SetString(reqBody.ClientPubkeyBN, 10)
	if !ok {
		global.Logger.Error("wrong clientPubkeyBN")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	serverPubkeyBN := new(big.Int)
	if len(reqBody.ServerSKSeed) > 0 {
		ServerSKSeedBN, ok := new(big.Int).SetString(reqBody.ServerSKSeed, 10)
		if !ok {
			global.Logger.Error("wrong ServerSKSeed")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		serverPubkeyBN, _ = p1.KeyGenRound1FromSeed(clientPubkeyBN, ServerSKSeedBN)
	} else {
		serverPubkeyBN, _ = p1.KeyGenRound1NoSeed(clientPubkeyBN)
	}

	resp := rest.P1KeygenRound1Response{
		ServerPubkeyBN: serverPubkeyBN.String(),
	}
	base_resp.JsonResponseSimple(c, resp)
}
