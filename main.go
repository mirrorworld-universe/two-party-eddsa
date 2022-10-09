package main

import (
	"main/internal/eddsa"
	"main/internal/p0"
	"main/internal/p1"
	"math/big"
)

func toLittleEdian(a []byte) (b []byte) {
	b = make([]byte, len(a))
	copy(b, a)
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
	}
	return b
}

func main() {
	//rnd, _ := new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	//b := rnd.Bytes() // big endian byte
	//reader := bytes.NewReader(b)
	//publicKey, privateKey, _ := eddsa.GenerateKey(reader)
	//fmt.Println(toLittleEdian(publicKey), privateKey)

	clientKeypair, keyAgg := p0.KeyGen()
	println("clientKeypair=", clientKeypair)
	//
	//println("\n\n************ SIGN now *************")
	//msg := "hello"
	//p0.Sign(&msg, clientKeypair, keyAgg)

	//R := "1dd9ad91d660e104fc02043e7bbe0c303f1bfc1c012689ab8c2d38c4ae6be0e7"
	//s := "7bf0d2eb8027a65988c43a4c79e70f3ab67eadf1a8a852b5cf34ef1ace192407"
	//pubkey := "790c23f4a2f065fa4cebf77a005f75ad7a528c8de4ca64e4e5c681c17663514e"
	//p0.Verify(&msg, &R, &s, &pubkey)
	bn, _ := new(big.Int).SetString("2D282E87852DE00E981EFAAD08A9F435CF41CDFD7BD9EB3DBFBE08AE5537160", 16)
	serverKeypair := eddsa.CreateKeyPairFromSeed(bn)
	println("serverKeypair=", serverKeypair.ToString())

	p1.Sign(serverKeypair, keyAgg)

}
