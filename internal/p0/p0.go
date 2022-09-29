package p0

import (
	"fmt"
	"main/internal/eddsa"
	"math/big"
)

func KeyGen() {
	fmt.Println("*************Client*************")
	rnd, _ := new(big.Int).SetString("5266194697103632731894445446481908111422432681065623019013231350200571873746", 10)
	clientKeypair := eddsa.CreateKeyPairFromSeed(rnd)
	clientPublicKey := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPublicKey)

	fmt.Println("*************Server*************")
	rnd, _ = new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	serverKeypair := eddsa.CreateKeyPairFromSeed(rnd)
	serverPublicKey := [32]byte{}
	serverKeypair.PublicKey.Ge.ToBytes(&serverPublicKey)

}
