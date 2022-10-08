package eddsa

import (
	"bytes"
	"math/big"
)

type ExpendedPrivateKey struct {
	Prefix     Ed25519Scalar
	PrivateKey Ed25519Scalar
}

type Keypair struct {
	PublicKey          Ed25519Point
	ExtendedPrivateKey ExpendedPrivateKey
}

func CreateKeyPairFromSeed(seed *big.Int) *Keypair {
	b := seed.Bytes() // big endian byte
	reader := bytes.NewReader(b)
	keypair, _ := GenerateKey(reader)
	return keypair
}

func (kp *Keypair) ToString() string {
	return "{\"PublicKey\"=" + kp.PublicKey.ToString() + ", \"ExtendedPrivateKey\": {\"Prefix\":" + kp.ExtendedPrivateKey.Prefix.ToString() + ", \"PrivateKey\": " + kp.ExtendedPrivateKey.PrivateKey.ToString()
}
