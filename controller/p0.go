package controller

import (
	"github.com/gin-gonic/gin"
	"main/internal/base_resp"
	"main/model/rest"
	"main/service/p0"
	"math/big"
	"net/http"
)

func P0KeyGenRound1(c *gin.Context) {
	clientKeypair, keyAgg := p0.KeyGenRound1NoSeed()

	clientPubkeyByte := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPubkeyByte)

	resp := rest.P0KeygenRound1Response{
		ClientPubkeyBN: new(big.Int).SetBytes(clientPubkeyByte[:]).String(),
		KeyAgg:         keyAgg.Apk.ToHexString(),
	}
	base_resp.JsonResponseSimple(c, resp)
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
