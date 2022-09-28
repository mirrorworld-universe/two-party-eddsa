package ed25519

import (
	"errors"
	"github.com/agl/ed25519/edwards25519"
)

type GeP3 = edwards25519.ExtendedGroupElement

type Ed25519Point struct {
	Purpose string
	Ge      edwards25519.ExtendedGroupElement
}

func GeP3FromBytesNegativeVartime(s *[32]byte) *GeP3 {
	y := new(edwards25519.FieldElement)
	edwards25519.FeFromBytes(y, s)

	z := new(edwards25519.FieldElement)
	edwards25519.FeOne(z)

	y_squared := FeSimpleSquare(y)
	u := FeSimpleSub(y_squared, z)
	v1 := FeSimpleMul(y_squared, FeD)
	v := FeSimpleAdd(v1, z)
	v_raise_3 := FeSimpleMul(FeSimpleSquare(v), v)
	v_raise_7 := FeSimpleMul(FeSimpleSquare(v_raise_3), v)
	uv7 := FeSimpleMul(v_raise_7, u)

	x := FeSimpleMul(FeSimpleMul(FePow25523(uv7), v_raise_3), u)
	vxx := FeSimpleMul(FeSimpleSquare(x), v)
	check := FeSimpleSub(vxx, u)
	if FeSimpleIsNonZero(check) {
		check2 := FeSimpleAdd(vxx, u)
		if FeSimpleIsNonZero(check2) {
			return nil
		}
		x = FeSimpleMul(x, FeSQRTM1)
	}

	if FeSimpleIsNegative(x) == ((s[31] >> 7) != 0) {
		x = FeSimpleNeg(x)
	}

	t := FeSimpleMul(x, y)
	p3 := GeP3{
		X: *x,
		Y: *y,
		Z: *z,
		T: *t,
	}
	return &p3
}

func ECPointGenerator() *Ed25519Point {
	vec_1 := [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var ge edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&ge, &vec_1)
	return &Ed25519Point{
		Purpose: "base_fe",
		Ge:      ge,
	}
}

func (e *Ed25519Point) ECPMul(fe *edwards25519.FieldElement) *Ed25519Point {
	vec0 := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	aScalarBytes := [32]byte{}
	edwards25519.FeToBytes(&aScalarBytes, fe)

	p2Point := new(edwards25519.ProjectiveGroupElement)
	edwards25519.GeDoubleScalarMultVartime(p2Point, &aScalarBytes, &e.Ge, &vec0)
	p2Bytes := [32]byte{}
	p2Point.ToBytes(&p2Bytes)
	p2Bytes[31] = p2Bytes[31] ^ (1 << 7)

	ge := GeP3FromBytesNegativeVartime(&p2Bytes)
	if ge == nil {
		panic(errors.New("GeP3FromBytesNegativeVartime check2 failed"))
	}

	return &Ed25519Point{
		Purpose: "scalar_point_mul",
		Ge:      *ge,
	}
}
