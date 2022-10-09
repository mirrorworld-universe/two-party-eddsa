package eddsa

import (
	"errors"
	"main/internal/agl_ed25519/edwards25519"
	"main/utils"
	"math/big"
)

type GeP3 = edwards25519.ExtendedGroupElement
type GeP1P1 = edwards25519.CompletedGroupElement
type GeP2 = edwards25519.ProjectiveGroupElement

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

func (e *Ed25519Point) ECPAdd(other *GeP3) *GeP1P1 {
	pkpk := new(GeP1P1)
	otherCached := new(edwards25519.CachedGroupElement)
	other.ToCached(otherCached)
	edwards25519.GeAdd(pkpk, &e.Ge, otherCached)
	return pkpk
}

func (e *Ed25519Point) ECPAddPoint(other *GeP3) *Ed25519Point {
	pkpk := e.ECPAdd(other)
	pkP2 := new(edwards25519.ProjectiveGroupElement)
	pkpk.ToProjective(pkP2)
	pkP2Bytes := [32]byte{}
	pkP2.ToBytes(&pkP2Bytes)
	pkP2Bytes[31] ^= 1 << 7
	ge := GeP3FromBytesNegativeVartime(&pkP2Bytes)
	if ge == nil {
		panic(errors.New("ECPAddPoint GeP3FromBytesNegativeVartime failed"))
	}

	p := &Ed25519Point{
		Purpose: "combine",
		Ge:      *ge,
	}
	return p
}

func ECPFromBytes(b *[32]byte) *Ed25519Point {
	geFromBytes := GeP3FromBytesNegativeVartime(b)
	geBytes := [32]byte{}
	geFromBytes.ToBytes(&geBytes)
	geFromBytes = GeP3FromBytesNegativeVartime(&geBytes)
	eight := ECSFromBigInt(new(big.Int).SetInt64(8))
	newPoint := Ed25519Point{
		Purpose: "random",
		Ge:      *geFromBytes,
	}
	newPoint2 := newPoint.ECPMul(&eight.Fe)
	return newPoint2
}

func (e *Ed25519Point) BytesCompressedToBigInt() *big.Int {
	bytes := &[32]byte{}
	e.Ge.ToBytes(bytes)
	return new(big.Int).SetBytes(bytes[:])
}

func (e *Ed25519Point) IsEqual(other *Ed25519Point) bool {
	aByte := [32]byte{}
	e.Ge.ToBytes(&aByte)

	bByte := [32]byte{}
	other.Ge.ToBytes(&bByte)

	aBN := new(big.Int).SetBytes(aByte[:])
	bBN := new(big.Int).SetBytes(bByte[:])
	return aBN.String() == bBN.String()
}

func (e *Ed25519Point) ToString() string {
	geBytes := [32]byte{}
	e.Ge.ToBytes(&geBytes)
	return "Purpose=" + e.Purpose + " ge.bytes=" + utils.BytesToStr(geBytes[:])
}
