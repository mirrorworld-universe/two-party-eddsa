package p0

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/eddsa"
	"main/internal/utils"
	"math/big"
)

//func tempLoadKey() (*eddsa.Keypair, *eddsa.KeyAgg) {
//
//	eightInverse := eddsa.ECSReverseBNToECS(new(big.Int).SetInt64(8))
//
//	rnd, _ := new(big.Int).SetString("5266194697103632731894445446481908111422432681065623019013231350200571873746", 10)
//	clientKeypair := eddsa.CreateKeyPairFromSeed(rnd)
//
//	aggPubKeyBytes := [32]byte{121, 12, 35, 244, 162, 240, 101, 250, 76, 235, 247, 122, 0, 95, 117, 173, 122, 82, 140, 141, 228, 202, 100, 228, 229, 198, 129, 193, 118, 99, 81, 78}
//	aggPubkey := eddsa.ECPFromBytes(&aggPubKeyBytes)
//	aggPubkey = aggPubkey.ECPMul(&eightInverse.Fe)
//	hashTBytes := [32]byte{233, 77, 67, 129, 29, 58, 37, 106, 36, 232, 48, 15, 76, 200, 132, 235, 167, 218, 242, 201, 195, 148, 83, 162, 158, 111, 87, 141, 120, 193, 14, 10}
//	hashFe := eddsa.ECSFromBigInt(new(big.Int).SetBytes(hashTBytes[:]))
//
//	keyAgg := eddsa.KeyAgg{
//		Apk:  *aggPubkey,
//		Hash: hashFe,
//	}
//
//	return clientKeypair, &keyAgg
//}

func Sign(msg *string, clientKeypair *eddsa.Keypair, keyAgg *eddsa.KeyAgg) {

	//clientKeypair, keyAgg := tempLoadKey()
	println("clientKeyPair=", clientKeypair.ToString(), " keyagg=", keyAgg.ToString())

	// round 1
	msgHash := sha256.Sum256([]byte(*msg))
	println("msgHash=", new(big.Int).SetBytes(msgHash[:]).String())

	clientEphemeralKey, clientSignFirstMsg, clientSignSecondMsg := eddsa.CreateEphemeralKeyAndCommit(clientKeypair, msgHash[:])
	println("clientEphemeralKey=", clientEphemeralKey.ToString(), ", clientSignFirstMsg=", clientSignFirstMsg.ToString()+", clientSignSecondMsg=", clientSignSecondMsg.ToString())

	// now send clientSignFirstMsg, msgHash, client public key to p1, and receive serverFirstSignMsg
	serverCommitment, _ := new(big.Int).SetString("84931746524459149992060349634228453990530694124359495037784716096273864068584", 10)
	serverSignFirstMsg := eddsa.SignFirstMsg{
		Commitment: *serverCommitment,
	}
	println("serverSignFirstMsg=", serverSignFirstMsg.ToString())

	// round 2
	// send clientSecondSignMsg to p1, get serverSignSecondMsg{R, blindFactor}
	eight := eddsa.ECSFromBigInt(new(big.Int).SetInt64(8))
	eightInverse := eight.ModInvert()
	serverSignSecondMsgRBytes := [32]byte{
		142, 144, 114, 134, 190, 107, 127, 90,
		212, 252, 156, 101, 121, 82, 106, 155,
		187, 60, 75, 220, 240, 209, 132, 217,
		100, 78, 252, 14, 20, 73, 153, 54,
	}
	serverSignSecondMsgR := eddsa.ECPFromBytes(&serverSignSecondMsgRBytes)
	serverSignSecondMsgR = serverSignSecondMsgR.ECPMul(&eightInverse.Fe)

	temp1 := [32]byte{
		169, 43, 89, 150, 255, 113, 182, 143,
		232, 177, 192, 27, 76, 61, 36, 72,
		121, 68, 213, 61, 241, 206, 20, 165,
		112, 33, 80, 6, 72, 206, 30, 83,
	}
	serverSignSecondMsgBF := new(big.Int).SetBytes(temp1[:])
	println("round2, server_sign_second_msg_R=", serverSignSecondMsgR.ToString(), ", serverSignSecondMsgBF=", serverSignSecondMsgBF.String())

	temp2 := [32]byte{
		29, 217, 173, 145, 214, 96, 225, 4,
		252, 2, 4, 62, 123, 190, 12, 48,
		63, 27, 252, 28, 1, 38, 137, 171,
		140, 45, 56, 196, 174, 107, 224, 231,
	}
	serverSigR := eddsa.ECPFromBytes(&temp2)
	serverSigR = serverSigR.ECPMul(&eightInverse.Fe)
	println("round2, serverSigR=", serverSigR.ToString())

	temp3 := [32]byte{
		124, 155, 253, 249, 189, 116, 9, 104,
		139, 154, 108, 227, 90, 150, 239, 201,
		172, 186, 250, 211, 86, 58, 200, 208,
		138, 102, 125, 137, 46, 247, 205, 10,
	}
	temp33 := utils.ReverseByteSlice(temp3[:])
	serverSigS := eddsa.ECSFromBigInt(new(big.Int).SetBytes(temp33))

	serverSignSecondMsg := eddsa.SignSecondMsg{
		R:           *serverSignSecondMsgR,
		BlindFactor: *serverSignSecondMsgBF,
	}
	serverSig := eddsa.Signature{
		R:      *serverSigR,
		SmallS: serverSigS,
	}
	println("round2, serverSignSecondMsg=", serverSignSecondMsg.ToString(), " serverSig=", serverSig.ToString())

	// check commiment
	isCommMatch := eddsa.CheckCommitment(
		&serverSignSecondMsg.R,
		&serverSignSecondMsg.BlindFactor,
		&serverSignFirstMsg.Commitment,
	)

	if !isCommMatch {
		panic(errors.New("commitment not match"))
	}

	// round 3
	ri := []eddsa.Ed25519Point{
		*serverSignSecondMsgR,
		clientSignSecondMsg.R,
	}
	rTot := eddsa.SigGetRTot(ri)
	println("rTot=", rTot.ToString())

	msgHash2 := msgHash[:]
	k := eddsa.SigK(rTot, &keyAgg.Apk, &msgHash2)
	println("k=", k.ToString())

	s2 := eddsa.PartialSign(
		&clientEphemeralKey.SmallR,
		clientKeypair,
		&k,
		&keyAgg.Hash,
		rTot,
	)
	println("s2=", s2.ToString())

	s := []eddsa.Signature{
		serverSig,
		s2,
	}
	sig := eddsa.AddSignatureParts(s)
	RBytes := [32]byte{}
	sig.R.Ge.ToBytes(&RBytes)
	sBytes := [32]byte{}
	edwards25519.FeToBytes(&sBytes, &sig.SmallS.Fe)
	println("sig=", sig.ToString(), " R: ", hex.EncodeToString(RBytes[:]), " s:", hex.EncodeToString(sBytes[:]))

	//clientPublicKeyBytes := [32]byte{}
	//clientKeypair.PublicKey.Ge.ToBytes(&clientPublicKeyBytes)
	//
	//temp := [][]byte{
	//	[]byte{2},
	//	utils.BigintToBytes32(&clientSignFirstMsg.Commitment),
	//	msgHash[:32],
	//	clientPublicKeyBytes[:],
	//}
	//
	//// send to p1
	//
	//buf := make([]byte, 32)
	//serverCommitmentBytes := []byte{}
	//serverCommitment := new(big.Int).SetBytes(serverCommitmentBytes)
	//serverSignFirstMsg := eddsa.SignFirstMsg{
	//	Commitment: *serverCommitment,
	//}
	//
	//// round 2
	//clientSignSecondMsgBytes := [32]byte{}
	//clientSignSecondMsg.R.Ge.ToBytes(&clientSignSecondMsgBytes)
	//temp = [][]byte{
	//	clientSignSecondMsgBytes[:],
	//	utils.BigintToBytes32(&clientSignSecondMsg.BlindFactor),
	//}
	//buf = utils.ConcatSlices(temp)
	// send to p1

}
