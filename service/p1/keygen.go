package p1

import (
	"main/global"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/eddsa"
	"main/model/db"
	"math/big"
)

func KeyGenRound1FromSeed(userId *string, clientPubkeyBN *big.Int, serverSKSeed *big.Int) (*big.Int, *eddsa.KeyAgg) {
	//println("P1 KeyGenRound1FromSeed=", serverSKSeed.String())
	return keyGenRound1Internal(userId, clientPubkeyBN, serverSKSeed)
}

func KeyGenRound1NoSeed(userId *string, clientPubkeyBN *big.Int) (*big.Int, *eddsa.KeyAgg) {
	//println("P1 KeyGenRound1NoSeed")
	ecsRndBytes := [32]byte{}
	edwards25519.FeToBytes(&ecsRndBytes, &eddsa.ECSNewRandom().Fe)
	serverSKSeed := new(big.Int).SetBytes(ecsRndBytes[:])

	return keyGenRound1Internal(userId, clientPubkeyBN, serverSKSeed)
}

func GenerateKeyAgg(clientPubkey *eddsa.Ed25519Point, serverPubkey *eddsa.Ed25519Point, idx uint8) *eddsa.KeyAgg {
	pks := []eddsa.Ed25519Point{
		*serverPubkey,
		*clientPubkey,
	}
	keyAgg := eddsa.KeyAggregationN(&pks, idx)
	return keyAgg
}

/**
Response: serverPubkey
DB: serverKeypair, aggPubkey
*/
func keyGenRound1Internal(userId *string, clientPubkeyBN *big.Int, serverSKSeed *big.Int) (*big.Int, *eddsa.KeyAgg) {
	clientPubkey := eddsa.NewECPSetFromBN(clientPubkeyBN)
	serverKeypair := eddsa.CreateKeyPairFromSeed(serverSKSeed)
	serverPubkeyByte := [32]byte{}
	serverKeypair.PublicKey.Ge.ToBytes(&serverPubkeyByte)
	serverPubkeyBN := new(big.Int).SetBytes(serverPubkeyByte[:])

	// start keyAgg
	keyAgg := GenerateKeyAgg(clientPubkey, &serverKeypair.PublicKey, global.PARTY_INDEX_P1)
	println("serverPubKeyBN: ", serverPubkeyBN.String(), " keyAgg=", keyAgg.ToString())

	// store serverSK, clientPubkey, keyAgg to db.
	wallet := db.MPCWallet{
		UserId:       *userId,
		PartyIdx:     1,
		SeedBN:       serverSKSeed.String(),
		KeyAggAPKBN:  keyAgg.Apk.ToBigInt().String(),
		KeyAggHashBN: keyAgg.Hash.ToBigInt().String(),
	}
	wallet.Create()

	return serverPubkeyBN, keyAgg
}
