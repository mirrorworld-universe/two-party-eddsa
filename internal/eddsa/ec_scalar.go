package eddsa

import (
	cryptorand "crypto/rand"
	"fmt"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/utils"
	"math/big"
)

const TWO_TIMES_SECRET_KEY_SIZE = 64

type Ed25519Scalar struct {
	Purpose string
	Fe      edwards25519.FieldElement
}

func ECSFromBigInt(n *big.Int) Ed25519Scalar {
	v := n.Bytes()
	if len(v) > TWO_TIMES_SECRET_KEY_SIZE {
		v = v[0:TWO_TIMES_SECRET_KEY_SIZE]
	}
	template := make([]byte, TWO_TIMES_SECRET_KEY_SIZE-len(v))
	template = append(template, v...)
	v = template
	v = utils.ReverseByteSlice(v)

	out := [32]byte{}
	v64 := [64]byte{}
	for i := 0; i < 64; i++ {
		v64[i] = v[i]
	}
	edwards25519.ScReduce(&out, &v64)
	fe := new(edwards25519.FieldElement)
	edwards25519.FeFromBytes(fe, &out)
	return Ed25519Scalar{
		Purpose: "from_big_int",
		Fe:      *fe,
	}
}

func (e *Ed25519Scalar) toBigInt() *big.Int {
	feBytes := [32]byte{}
	edwards25519.FeToBytes(&feBytes, &e.Fe)

	// reverse fe_bytes
	for i, j := 0, len(feBytes)-1; i < j; i, j = i+1, j-1 {
		feBytes[i], feBytes[j] = feBytes[j], feBytes[i]
	}

	ret := new(big.Int).SetBytes(feBytes[:])
	return ret
}

func q() *big.Int {
	qBytesArray := [32]byte{237, 211, 245, 92, 26, 99, 18, 88, 214, 156, 247, 162, 222, 249, 222, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16}
	lFe := new(edwards25519.FieldElement)
	edwards25519.FeFromBytes(lFe, &qBytesArray)
	lFeScalar := Ed25519Scalar{
		Purpose: "q",
		Fe:      *lFe,
	}
	return lFeScalar.toBigInt()
}

//func from(n *big.Int) *Ed25519Scalar {
//	n_bytes := n.Bytes()
//	n_bytes_64 = n_bytes[:]
//	n_bytes_r := utils.ReverseByteSlice(n_bytes)
//	out := [32]byte{}
//	edwards25519.ScReduce(&out, &n_bytes_r)
//}

func ECSNewRandom() *Ed25519Scalar {
	// sample_below()
	reader := cryptorand.Reader
	rndBn, _ := cryptorand.Int(reader, q())
	bn8 := big.NewInt(8)
	rndBnMul := new(big.Int).Mul(rndBn, bn8)
	rndBnMul8 := new(big.Int).Mod(rndBnMul, q())
	ret := ECSFromBigInt(rndBnMul8)
	return &ret
}

func ECSReverseBNToECS(bn *big.Int) Ed25519Scalar {
	intBytes := bn.Bytes()
	intBytes = utils.ReverseByteSlice(intBytes)
	scalarOut := new(big.Int).SetBytes(intBytes)
	ret := ECSFromBigInt(scalarOut)
	return ret
}

func (e *Ed25519Scalar) Print() {
	feBytes := [32]byte{}
	edwards25519.FeToBytes(&feBytes, &e.Fe)
	fmt.Println("Purpose=", e.Purpose, " fe.bytes=", feBytes)
}

func (e *Ed25519Scalar) ToString() string {
	feBytes := [32]byte{}
	edwards25519.FeToBytes(&feBytes, &e.Fe)
	return "Purpose=" + e.Purpose + " fe.bytes=" + utils.BytesToStr(feBytes[:])
}
