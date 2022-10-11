package p1

import (
	"errors"
	"main/internal/eddsa"
	"math/big"
)

func SignRound1(serverKeypair *eddsa.Keypair, msgHash *big.Int) (*eddsa.EphemeralKey, *eddsa.SignFirstMsg) {
	serverEphemeralKey, serverSignFirstMsg, _ := eddsa.CreateEphemeralKeyAndCommit(serverKeypair, msgHash.Bytes())
	return &serverEphemeralKey, &serverSignFirstMsg
}

func SignRound2(
	clientCommitment *big.Int,
	serverKeypair *eddsa.Keypair,
	msgHash *big.Int,
	clientSignSecondMsgR *eddsa.Ed25519Point,
	clientSignSecondMsgBF *big.Int,
	keyAgg *eddsa.KeyAgg,
) (*eddsa.SignSecondMsg, *eddsa.Signature) {

	//// check commiment
	isCommMatch := eddsa.CheckCommitment(
		clientSignSecondMsgR,
		clientSignSecondMsgBF,
		clientCommitment,
	)
	if !isCommMatch {
		panic(errors.New("commitment not match"))
	}

	serverEphemeralKey, _, serverSignSecondMsg := eddsa.CreateEphemeralKeyAndCommit(serverKeypair, msgHash.Bytes())
	ri := []eddsa.Ed25519Point{
		serverSignSecondMsg.R,
		*clientSignSecondMsgR,
	}
	rTot := eddsa.SigGetRTot(ri)
	println("rTot=", rTot.ToString())

	msgHashBytes := msgHash.Bytes()
	k := eddsa.SigK(rTot, &keyAgg.Apk, &msgHashBytes)
	println("k=", k.ToString())

	s1 := eddsa.PartialSign(
		&serverEphemeralKey.SmallR,
		serverKeypair,
		&k,
		&keyAgg.Hash,
		rTot,
	)
	println("rTot=", rTot.ToString(), "k=", k.ToString(), "s1=", s1.ToString())
	return &serverSignSecondMsg, &s1
}

func Sign(serverKeypair *eddsa.Keypair, keyAgg *eddsa.KeyAgg) {
	//clientKeypair, keyAgg := tempLoadKey()
	println("serverKeypair=", serverKeypair.ToString(), " keyAgg=", keyAgg.ToString())

	temp1 := [32]byte{
		196, 206, 142, 211, 33, 209, 253, 148,
		123, 34, 148, 86, 59, 179, 45, 242,
		225, 141, 253, 97, 58, 156, 84, 142,
		123, 77, 96, 1, 160, 16, 53, 164,
	}
	clientCommitment := new(big.Int).SetBytes(temp1[:])
	temp2 := [32]byte{
		44, 242, 77, 186, 95, 176, 163, 14,
		38, 232, 59, 42, 197, 185, 226, 158,
		27, 22, 30, 92, 31, 167, 66, 94,
		115, 4, 51, 98, 147, 139, 152, 36,
	}
	msgHash := new(big.Int).SetBytes(temp2[:])
	println("msgHash=", msgHash.String(), " clientCommitment=", clientCommitment.String())

	eight := eddsa.ECSFromBigInt(new(big.Int).SetInt64(8))
	eightInverse := eight.ModInvert()

	temp3 := [32]byte{
		55, 94, 156, 213, 171, 71, 43, 180,
		25, 210, 117, 204, 69, 176, 139, 18,
		223, 232, 74, 151, 45, 31, 77, 169,
		236, 104, 205, 15, 156, 87, 239, 134,
	}

	clientPubkeyBN := new(big.Int).SetBytes(temp3[:])
	clientPubkey := eddsa.NewECPSetFromBN(clientPubkeyBN)
	clientPubkey = clientPubkey.ECPMul(&eightInverse.Fe)
	println("clientPubkey=", clientPubkey.ToString(), ", clientPubkeyBN=", clientPubkeyBN.String())

	// round 1

	serverEphemeralKey, serverSignFirstMsg, serverSignSecondMsg := eddsa.CreateEphemeralKeyAndCommit(serverKeypair, msgHash.Bytes())
	println("serverEphemeralKey=", serverEphemeralKey.ToString(), ", serverSignFirstMsg=", serverSignFirstMsg.ToString()+", serverSignSecondMsg=", serverSignSecondMsg.ToString())

	// now send clientSignFirstMsg, msgHash, client public key to p1, and receive serverFirstSignMsg

	// round 2

	clientSignSecondMsgRBytes := [32]byte{
		157, 128, 104, 28, 9, 210, 124, 222,
		144, 57, 146, 136, 200, 8, 228, 125,
		119, 216, 89, 255, 247, 124, 203, 150,
		67, 246, 66, 228, 108, 121, 46, 246,
	}
	clientSignSecondMsgR := eddsa.ECPFromBytes(&clientSignSecondMsgRBytes)
	clientSignSecondMsgR = clientSignSecondMsgR.ECPMul(&eightInverse.Fe)

	temp1 = [32]byte{
		169, 43, 89, 150, 255, 113, 182, 143,
		232, 177, 192, 27, 76, 61, 36, 72,
		121, 68, 213, 61, 241, 206, 20, 165,
		112, 33, 80, 6, 72, 206, 30, 83,
	}
	clientSignSecondMsgBF := new(big.Int).SetBytes(temp1[:])
	println("round2, clientSignSecondMsgR=", clientSignSecondMsgR.ToString(), ", clientSignSecondMsgBF=", clientSignSecondMsgBF.String())

	// check commiment
	isCommMatch := eddsa.CheckCommitment(
		clientSignSecondMsgR,
		clientSignSecondMsgBF,
		clientCommitment,
	)

	if !isCommMatch {
		panic(errors.New("commitment not match"))
	}

	ri := []eddsa.Ed25519Point{
		serverSignSecondMsg.R,
		*clientSignSecondMsgR,
	}
	rTot := eddsa.SigGetRTot(ri)
	println("rTot=", rTot.ToString())

	msgHashBytes := msgHash.Bytes()
	k := eddsa.SigK(rTot, &keyAgg.Apk, &msgHashBytes)
	println("k=", k.ToString())

	s1 := eddsa.PartialSign(
		&serverEphemeralKey.SmallR,
		serverKeypair,
		&k,
		&keyAgg.Hash,
		rTot,
	)
	println("rTot=", rTot.ToString(), "k=", k.ToString(), "s1=", s1.ToString())

	// send to client
}
