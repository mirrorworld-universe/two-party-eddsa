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

type EphemeralKey struct {
	R      Ed25519Point
	SmallR Ed25519Scalar
}

type SignFirstMsg struct {
	Commitment big.Int
}

type SignSecondMsg struct {
	R           Ed25519Point
	BlindFactor big.Int
}

func CreateEphemeralKeyAndCommit(key *Keypair, message []byte) (EphemeralKey, SignFirstMsg, SignSecondMsg) {

	prefixBytes := [32]byte{}
	edwards25519.FeToBytes(&prefixBytes, &key.ExtendedPrivateKey.Prefix.Fe)
	ecsRndBytes := [32]byte{}
	edwards25519.FeToBytes(&ecsRndBytes, &ECSNewRandom().Fe)
	bytes := [][]byte{
		new(big.Int).SetInt64(2).Bytes(),
		prefixBytes[:],
		message,
		ecsRndBytes[:],
	}
	concatBytes := utils.ConcatSlices(bytes)
	r := sha512.Sum512(concatBytes)
	rInt := new(big.Int).SetBytes(r[:])
	r2 := ECSReverseBNToECS(rInt)

	ecPoint := ECPointGenerator()
	R := ecPoint.ECPMul(&r2.Fe)

	hashCom := CreateCommitment(R.BytesCompressedToBigInt())

	return EphemeralKey{
			SmallR: r2,
			R:      *R,
		}, SignFirstMsg{
			Commitment: hashCom.Commitment,
		}, SignSecondMsg{
			R:           *R,
			BlindFactor: hashCom.BlindFactor,
		}
}

func SigGetRTot(R []Ed25519Point) *Ed25519Point {
	sum := new(Ed25519Point)
	for _, v := range R {
		sum = sum.ECPAddPoint(&v.Ge)
	}
	return sum
}

func SigK(R_tot *Ed25519Point, apk *Ed25519Point, message *[]byte) Ed25519Scalar {
	temp := [][]byte{
		R_tot.BytesCompressedToBigInt().Bytes(),
		apk.BytesCompressedToBigInt().Bytes(),
	}
	k := sha512.Sum512(utils.ConcatSlices(temp))
	kBN := new(big.Int).SetBytes(k[:])
	k2 := ECSReverseBNToECS(kBN)
	return k2
}

func PartialSign(r *Ed25519Scalar, key *Keypair, k *Ed25519Scalar, a *Ed25519Scalar, R_tot *Ed25519Point) Signature {
	kMulSk := k.Mul(&key.ExtendedPrivateKey.PrivateKey)
	kMulSkMulAi := kMulSk.Mul(a)
	s := r.Add(&kMulSkMulAi)
	return Signature{
		R:      *R_tot,
		SmallS: s,
	}
}

func AddSignatureParts(sigs []Signature) Signature {
	candidateR := sigs[0].R
	for _, x := range sigs {
		if x.R != candidateR {
			panic("R not equal")
		}
	}
	sum := new(Ed25519Scalar)
	for _, v := range sigs {
		temp := sum.Add(&v.SmallS)
		sum = &temp
	}
	return Signature{
		SmallS: *sum,
		R:      sigs[0].R,
	}
}
