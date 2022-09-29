package eddsa

import (
	"fmt"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/utils"
	"math/big"
)

type Ed25519Scalar struct {
	Purpose string
	Fe      edwards25519.FieldElement
}

func ECSFromBigInt(n *big.Int) Ed25519Scalar {
	nBytes := n.Bytes()
	nBytes = utils.ReverseByteSlice(nBytes)
	var nBytes64 [64]byte
	for i := 0; i < len(nBytes); i++ {
		nBytes64[i] = nBytes[i]
	}

	out := [32]byte{}
	edwards25519.ScReduce(&out, &nBytes64)
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
