package p0

import (
	"crypto/sha256"
	"main/internal/eddsa"
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
