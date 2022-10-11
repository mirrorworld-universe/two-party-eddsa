package controller

import (
	"github.com/gin-gonic/gin"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/base_resp"
	"main/internal/binding"
	"main/internal/eddsa"
	"main/model/rest"
	"main/service/p1"
	"main/utils"
	"math/big"
)

func P1KeyGenRound1(c *gin.Context) {
	reqBody := rest.P1KeygenRound1Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}
	clientPubkeyBN, _ := new(big.Int).SetString(reqBody.ClientPubkeyBN, 10)
	serverPubkeyBN := new(big.Int)
	var keyAgg *eddsa.KeyAgg

	if len(reqBody.ServerSKSeed) > 0 {
		ServerSKSeedBN, _ := new(big.Int).SetString(reqBody.ServerSKSeed, 10)
		serverPubkeyBN, keyAgg = p1.KeyGenRound1FromSeed(clientPubkeyBN, ServerSKSeedBN)
	} else {
		serverPubkeyBN, keyAgg = p1.KeyGenRound1NoSeed(clientPubkeyBN)
	}
	println("server keyAgg=", keyAgg.ToString())
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
	serverEphemeralKey, serverSignFirstMsg := p1.SignRound1(serverKeypair, msgHash)
	println("serverEphemeralKey=", serverEphemeralKey.ToString(), "serverSignFirstMsg=", serverSignFirstMsg.ToString())

	resp := rest.P1SignRound1Response{
		ServerSignFirstMsgCommitmentBN: serverSignFirstMsg.Commitment.String(),
	}
	base_resp.JsonResponseSimple(c, resp)
}

func P1SignRound2(c *gin.Context) {
	eight := eddsa.ECSFromBigInt(new(big.Int).SetInt64(8))
	eightInverse := eight.ModInvert()

	reqBody := rest.P1SignRound2Req{}
	if err := binding.BindJson(c, &reqBody); err != nil {
		return
	}

	msgHash, _ := new(big.Int).SetString(reqBody.MsgHashBN, 10)
	clientCommitment, _ := new(big.Int).SetString(reqBody.ClientSignFirstMsgCommitmentBN, 10)
	clientSignSecondMsgRBN, _ := new(big.Int).SetString(reqBody.ClientSignSecondMsgRBN, 10)
	clientSignSecondMsgRByte := clientSignSecondMsgRBN.Bytes()
	clientSignSecondMsgR := eddsa.ECPFromBytes((*[32]byte)(clientSignSecondMsgRByte))
	clientSignSecondMsgR = clientSignSecondMsgR.ECPMul(&eightInverse.Fe)

	clientSignSecondMsgBF, _ := new(big.Int).SetString(reqBody.ClientSignSecondMsgBF32BN, 10)

	// hardcode
	// we should read serverSKSeed, keyAgg from db.
	ServerSKSeedBN, _ := new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	serverKeypair := eddsa.CreateKeyPairFromSeed(ServerSKSeedBN)
	clientPubkeyBN, _ := new(big.Int).SetString(reqBody.ClientPubkeyBN, 10)
	_, keyAgg := p1.KeyGenRound1FromSeed(clientPubkeyBN, ServerSKSeedBN)

	serverSignSecondMsg, s1 := p1.SignRound2(
		clientCommitment,
		serverKeypair,
		msgHash,
		clientSignSecondMsgR,
		clientSignSecondMsgBF,
		keyAgg,
	)

	temp32 := [32]byte{}
	serverSignSecondMsg.R.Ge.ToBytes(&temp32)
	serverSignSecondMsgR := new(big.Int).SetBytes(temp32[:]).String()

	bf32Byte := utils.BigintToBytes32(&serverSignSecondMsg.BlindFactor)

	s1.R.Ge.ToBytes(&temp32)
	serverSigRBN := new(big.Int).SetBytes(temp32[:]).String()

	edwards25519.FeToBytes(&temp32, &s1.SmallS.Fe)
	ServerSigSmallSBN := new(big.Int).SetBytes(temp32[:]).String()

	resp := rest.P1SignRound2Response{
		ServerSignSecondMsgR:    serverSignSecondMsgR,
		ServerSignSecondMsgBF32: new(big.Int).SetBytes(bf32Byte[:]).String(),
		ServerSigRBN:            serverSigRBN,
		ServerSigSmallSBN:       ServerSigSmallSBN,
	}
	base_resp.JsonResponseSimple(c, resp)
}
