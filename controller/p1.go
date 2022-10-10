package controller

import (
	"github.com/gin-gonic/gin"
	"main/global"
	"main/internal/base_resp"
	"main/internal/binding"
	"main/internal/eddsa"
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

func P1SignRound1(c *gin.Context) {
	reqBody := rest.P1SignRound1Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}
	clientCommitment, ok := new(big.Int).SetString(reqBody.ClientSignFirstMsgCommitmentBN, 10)
	println("ClientSignFirstMsgCommitmentBN=", clientCommitment.String())
	if !ok {
		global.Logger.Error("wrong ClientSignFirstMsgCommitmentBN")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	msgHash, ok := new(big.Int).SetString(reqBody.MsgHashBN, 10)
	if !ok {
		global.Logger.Error("wrong MsgHashBN")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	clientPubkeyBN, ok := new(big.Int).SetString(reqBody.ClientPubkeyBN, 10)
	if !ok {
		global.Logger.Error("wrong clientPubkeyBN")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	clientPubkey := eddsa.NewECPSetFromBN(clientPubkeyBN)
	println("clientPubkey=", clientPubkey.ToString())

	// hardcode
	ServerSKSeedBN, ok := new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	serverKeypair := eddsa.CreateKeyPairFromSeed(ServerSKSeedBN)
	serverEphemeralKey, serverSignFirstMsg, serverSignSecondMsg := eddsa.CreateEphemeralKeyAndCommit(serverKeypair, msgHash.Bytes())
	println("serverEphemeralKey=", serverEphemeralKey.ToString(), "serverSignFirstMsg=", serverSignFirstMsg.ToString(), "serverSignSecondMsg=", serverSignSecondMsg.ToString())

	resp := rest.P1SignRound1Response{
		ServerSignFirstMsgCommitmentBN: serverSignFirstMsg.Commitment.String(),
	}
	base_resp.JsonResponseSimple(c, resp)
}
