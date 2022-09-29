package eddsa

import (
	"crypto/sha512"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/utils"
	"math/big"
)

type Signature struct {
	R      Ed25519Point
	SmallS Ed25519Scalar
}

func CreateEphemeralKeyAndCommit(key *Keypair, message *[]byte) {

	prefixBytes := [32]byte{}
	edwards25519.FeToBytes(&prefixBytes, &key.ExtendedPrivateKey.Prefix.Fe)
	ecsRndBytes := [32]byte{}
	edwards25519.FeToBytes(&ecsRndBytes, &ECSNewRandom().Fe)
	bytes := [][]byte{
		new(big.Int).SetInt64(2).Bytes(),
		prefixBytes[:],
		*message,
		ecsRndBytes[:],
	}
	concatBytes := utils.ConcatSlices(bytes)
	r := sha512.Sum512(concatBytes)
	rInt := new(big.Int).SetBytes(r[:])
	r2 := ECSReverseBNToECS(rInt)

	ecPoint := ECPointGenerator()
	R := ecPoint.ECPMul(&r2.Fe)
}
