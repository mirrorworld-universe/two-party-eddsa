package eddsa

import "math/big"

type KeyAgg struct {
	Apk  Ed25519Point
	Hash Ed25519Scalar
}

func NewKeyAggFromBNs(apkBN *big.Int, hashBN *big.Int) *KeyAgg {
	keyAgg := KeyAgg{
		Apk:  *NewECPSetFromBN(apkBN),
		Hash: *NewECSSetFromBN(hashBN),
	}
	return &keyAgg
}

func (ka *KeyAgg) ToString() string {
	return "{" +
		"\"Apk\": " + ka.Apk.ToString() + ", " +
		"\"Hash\": " + ka.Hash.ToString() +
		"}"
}
