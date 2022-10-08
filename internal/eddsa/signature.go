package eddsa

import (
	"crypto/sha512"
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

	prefixBN := key.ExtendedPrivateKey.Prefix.ToBigInt()
	rnd, _ := new(big.Int).SetString("2030282828107764592039879086147438423373605693185722406299485015238703754456", 10)
	//ecsRndBytes := [32]byte{}
	//edwards25519.FeToBytes(&ecsRndBytes, &ECSNewRandom().Fe)
	bytes := [][]byte{
		new(big.Int).SetInt64(2).Bytes(),
		prefixBN.Bytes(),
		message,
		//ecsRndBytes[:],
		rnd.Bytes(),
	}
	println("1=", new(big.Int).SetInt64(2).String(),
		"2=", prefixBN.String(),
		"3=", new(big.Int).SetBytes(message).String(),
		"4=", new(big.Int).SetBytes(rnd.Bytes()).String(),
	)
	concatBytes := utils.ConcatSlices(bytes)
	r := sha512.Sum512(concatBytes)
	rInt := new(big.Int).SetBytes(r[:])
	r2 := ECSReverseBNToECS(rInt)
	println("CreateEphemeralKeyAndCommit, r=", new(big.Int).SetBytes(r[:]).String(),
		"r2=", r2.ToString(),
	)

	ecPoint := ECPointGenerator()
	R := ecPoint.ECPMul(&r2.Fe)

	//hashCom := CreateCommitment(R.BytesCompressedToBigInt())
	// hashcode
	blindFactor, _ := new(big.Int).SetString("76517464160675767839318574288422328116452541159689926027818280551122440429139", 10)
	commitment := CreateCommitmentWithUserDefinedRandomness(R.BytesCompressedToBigInt(), blindFactor)

	return EphemeralKey{
			SmallR: r2,
			R:      *R,
		}, SignFirstMsg{
			Commitment: *commitment,
		}, SignSecondMsg{
			R:           *R,
			BlindFactor: *blindFactor,
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

func (k *EphemeralKey) ToString() string {
	return "{\"R\": " + k.R.ToString() + "," +
		"\"SmallR\": " + k.SmallR.ToString() +
		"}"
}

func (m *Signature) ToString() string {
	return "{\"R\":" + m.R.ToString() + "," +
		"\"SmallS\":" + m.SmallS.ToString() +
		"}"
}

func (m *SignFirstMsg) ToString() string {
	return "{\"Commitment\":" + m.Commitment.String() + "}"
}

func (m *SignSecondMsg) ToString() string {
	return "{\"R\":" + m.R.ToString() + "," +
		"\"BlindFactor\":" + m.BlindFactor.String() +
		"}"
}
