package p0

import (
	"encoding/hex"
	"fmt"
	"main/internal/eddsa"
	"math/big"
)

func KeyGen() (*eddsa.Keypair, *eddsa.KeyAgg) {
	fmt.Println("*************Client*************")
	rnd, _ := new(big.Int).SetString("5266194697103632731894445446481908111422432681065623019013231350200571873746", 10)
	clientKeypair := eddsa.CreateKeyPairFromSeed(rnd)
	clientPublicKeyBytes := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPublicKeyBytes)

	fmt.Println("*************Server*************")
	rnd, _ = new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	serverKeypair := eddsa.CreateKeyPairFromSeed(rnd)
	serverPublicKeyBytes := [32]byte{}
	serverKeypair.PublicKey.Ge.ToBytes(&serverPublicKeyBytes)

	// start aggregate
	pks := []eddsa.Ed25519Point{
		serverKeypair.PublicKey, // partyIdx=0
		clientKeypair.PublicKey, // partyIdx=1
	}
	keyAgg := eddsa.KeyAggregationN(&pks, 1)
	aggPubKeyBytes := [32]byte{}
	keyAgg.Apk.Ge.ToBytes(&aggPubKeyBytes)
	fmt.Println("aggregated_pukey=", hex.EncodeToString(aggPubKeyBytes[:]))
	fmt.Println("key_agg=", keyAgg.ToString())

	// @TODO save to db
	return clientKeypair, keyAgg

	// continue sign process
}
