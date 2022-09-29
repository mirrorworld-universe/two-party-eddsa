package eddsa

import (
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
