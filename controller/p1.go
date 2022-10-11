package controller

import (
	"github.com/gin-gonic/gin"
	"main/internal/base_resp"
	"main/internal/binding"
	"main/internal/eddsa"
	"main/model/rest"
	"main/service/p1"
	"math/big"
)

func P1KeyGenRound1(c *gin.Context) {
	reqBody := rest.P1KeygenRound1Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}
	clientPubkeyBN, _ := new(big.Int).SetString(reqBody.ClientPubkeyBN, 10)
	serverPubkeyBN := new(big.Int)
	if len(reqBody.ServerSKSeed) > 0 {
		ServerSKSeedBN, _ := new(big.Int).SetString(reqBody.ServerSKSeed, 10)
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
	var reqBody rest.P1SignRound1Req
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}

	msgHash, _ := new(big.Int).SetString(reqBody.MsgHashBN, 10)
	clientPubkeyBN, _ := new(big.Int).SetString(reqBody.ClientPubkeyBN, 10)
	clientPubkey := eddsa.NewECPSetFromBN(clientPubkeyBN)
	println("clientPubkey=", clientPubkey.ToString())

	// hardcode
	ServerSKSeedBN, _ := new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	serverKeypair := eddsa.CreateKeyPairFromSeed(ServerSKSeedBN)
	serverEphemeralKey, serverSignFirstMsg, serverSignSecondMsg := p1.SignRound1(serverKeypair, msgHash)
	println("serverEphemeralKey=", serverEphemeralKey.ToString(), "serverSignFirstMsg=", serverSignFirstMsg.ToString(), "serverSignSecondMsg=", serverSignSecondMsg.ToString())

	resp := rest.P1SignRound1Response{
		ServerSignFirstMsgCommitmentBN: serverSignFirstMsg.Commitment.String(),
	}
	base_resp.JsonResponseSimple(c, resp)
}

func P1SignRound2(c *gin.Context) {
	//reqBody := rest.P1SignRound2Req{}
	//clientCommitment, _ := new(big.Int).SetString(reqBody.ClientSignFirstMsgCommitmentBN, 10)

}
