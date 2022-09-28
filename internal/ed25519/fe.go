package ed25519

import "github.com/agl/ed25519/edwards25519"

/**
Utils for FieldElement
*/

var FeD = &edwards25519.FieldElement{-10913610, 13857413, -15372611, 6949391, 114729, -8787816, -6275908, -3247719, -18696448, -12055116}
var FeSQRTM1 = &edwards25519.FieldElement{-32595792, -7943725, 9377950, 3500415, 12389472, -272473, -25146209, -2005654, 326686, 11406482}

func FeSquareN(fe *edwards25519.FieldElement, n int) *edwards25519.FieldElement {
	//  = 2 ** n
	now := new(edwards25519.FieldElement)
	edwards25519.FeCopy(now, fe)
	for i := 0; i < n; i++ {
		temp := new(edwards25519.FieldElement)
		temp2 := new(edwards25519.FieldElement)
		edwards25519.FeCopy(temp, now)
		edwards25519.FeSquare(temp2, temp)
		edwards25519.FeCopy(now, temp2)
	}
	return now
}

func FeSimpleIsNegative(fe *edwards25519.FieldElement) bool {
	return edwards25519.FeIsNegative(fe) != 0
}

func FeSimpleNeg(fe *edwards25519.FieldElement) *edwards25519.FieldElement {
	a := new(edwards25519.FieldElement)
	edwards25519.FeNeg(a, fe)
	return a
}

func FeSimpleMul(a *edwards25519.FieldElement, b *edwards25519.FieldElement) *edwards25519.FieldElement {
	c := new(edwards25519.FieldElement)
	edwards25519.FeMul(c, a, b)
	return c
}

func FeSimpleSquare(a *edwards25519.FieldElement) *edwards25519.FieldElement {
	b := new(edwards25519.FieldElement)
	edwards25519.FeSquare(b, a)
	return b
}

func FeSimpleSub(a *edwards25519.FieldElement, b *edwards25519.FieldElement) *edwards25519.FieldElement {
	c := new(edwards25519.FieldElement)
	edwards25519.FeSub(c, a, b)
	return c
}

func FeSimpleIsNonZero(a *edwards25519.FieldElement) bool {
	return edwards25519.FeIsNonZero(a) == 1
}

func FeSimpleAdd(a *edwards25519.FieldElement, b *edwards25519.FieldElement) *edwards25519.FieldElement {
	c := new(edwards25519.FieldElement)
	edwards25519.FeAdd(c, a, b)
	return c
}

func FePow25523(fe *edwards25519.FieldElement) *edwards25519.FieldElement {
	z2 := FeSimpleSquare(fe)
	z8 := FeSquareN(z2, 2)
	z9 := FeSimpleMul(fe, z8)
	z11 := FeSimpleMul(z9, z2)
	z22 := FeSimpleSquare(z11)
	z_5_0 := FeSimpleMul(z9, z22)
	z_10_5 := FeSquareN(z_5_0, 5)
	z_10_0 := FeSimpleMul(z_10_5, z_5_0)

	z_20_10 := FeSquareN(z_10_0, 10)
	z_20_0 := FeSimpleMul(z_20_10, z_10_0)

	z_40_20 := FeSquareN(z_20_0, 20)
	z_40_0 := FeSimpleMul(z_40_20, z_20_0)

	z_50_10 := FeSquareN(z_40_0, 10)
	z_50_0 := FeSimpleMul(z_50_10, z_10_0)

	z_100_50 := FeSquareN(z_50_0, 50)
	z_100_0 := FeSimpleMul(z_100_50, z_50_0)

	z_200_100 := FeSquareN(z_100_0, 100)
	z_200_0 := FeSimpleMul(z_200_100, z_100_0)

	z_250_50 := FeSquareN(z_200_0, 50)
	z_250_0 := FeSimpleMul(z_250_50, z_50_0)

	z_252_2 := FeSquareN(z_250_0, 2)
	z_252_3 := FeSimpleMul(z_252_2, fe)

	return z_252_3
}
