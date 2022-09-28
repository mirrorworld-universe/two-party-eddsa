package main

import (
	"bytes"
	"fmt"
	"main/internal/ed25519"
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
	rnd, _ := new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	b := rnd.Bytes() // big endian byte
	reader := bytes.NewReader(b)
	publicKey, privateKey, _ := ed25519.GenerateKey(reader)
	fmt.Println(toLittleEdian(publicKey), privateKey)
}
