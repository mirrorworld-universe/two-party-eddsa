package p1

import (
	"main/global"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/eddsa"
	"math/big"
)

func KeyGenRound1FromSeed(clientPubkeyBN *big.Int, serverSKSeed *big.Int) (*big.Int, *eddsa.KeyAgg) {
	//println("P1 KeyGenRound1FromSeed=", serverSKSeed.String())
	return keyGenRound1Internal(clientPubkeyBN, serverSKSeed)
}

func KeyGenRound1NoSeed(clientPubkeyBN *big.Int) (*big.Int, *eddsa.KeyAgg) {
	//println("P1 KeyGenRound1NoSeed")
	ecsRndBytes := [32]byte{}
	edwards25519.FeToBytes(&ecsRndBytes, &eddsa.ECSNewRandom().Fe)
	serverSKSeed := new(big.Int).SetBytes(ecsRndBytes[:])

	return keyGenRound1Internal(clientPubkeyBN, serverSKSeed)
}

/**
Response: serverPubkey
DB: serverKeypair, aggPubkey
*/
func keyGenRound1Internal(clientPubkeyBN *big.Int, serverSKSeed *big.Int) (*big.Int, *eddsa.KeyAgg) {
	clientPubkey := eddsa.NewECPSetFromBN(clientPubkeyBN)
	//println("clientPubkey from bn:", clientPubkey.ToString())
	//eight := eddsa.ECSFromBigInt(new(big.Int).SetInt64(8))
	//eightInverse := eight.ModInvert()
	//clientPubkey = clientPubkey.ECPMul(&eightInverse.Fe)

	serverKeypair := eddsa.CreateKeyPairFromSeed(serverSKSeed)
	serverPubkeyByte := [32]byte{}
	serverKeypair.PublicKey.Ge.ToBytes(&serverPubkeyByte)
	serverPubkeyBN := new(big.Int).SetBytes(serverPubkeyByte[:])

	// start keyAgg
	pks := []eddsa.Ed25519Point{
		serverKeypair.PublicKey,
		*clientPubkey,
	}
	keyAgg := eddsa.KeyAggregationN(&pks, global.PARTY_INDEX_P1)
	println("serverPubKeyBN: ", serverPubkeyBN.String(), " keyAgg=", keyAgg.ToString())

	return serverPubkeyBN, keyAgg
}
