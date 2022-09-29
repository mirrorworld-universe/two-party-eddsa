package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func Reverse32Slice(a *[32]byte) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func ReverseByteSlice(a []byte) []byte {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func BigIntSliceToString(s *[]big.Int) string {
	ret := "["
	for _, v := range *s {
		ret += v.String() + ","
	}
	ret += "]"
	return ret
}

func BytesToInts(s []byte) []int {
	r := make([]int, len(s))
	for i, v := range s {
		r[i] = int(v)
	}
	return r
}

func Int2str(i int) string {
	return strconv.Itoa(i)
}

func BytesToStr(s []byte) string {
	i := BytesToInts(s)
	r := "["
	for _, v := range i {
		r += Int2str(v) + ", "
	}
	return r + "]"
}

func StringToBytes(s *string) []byte {
	return []byte(*s)
}

func ConcatSlices(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func BigintSample(bitSize uint) *big.Int {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	n, _ := rand.Int(rand.Reader, max)

	bytes := (bitSize-1)/8 + 1
	buf := make([]byte, bytes)
	n.FillBytes(buf)

	n2 := new(big.Int).SetBytes(buf)
	n2.Rsh(n2, bytes*8-bitSize)
	return n2
}

func BigintToBytes32(n *big.Int) []byte {
	buf := []byte{}
	bBytes := n.Bytes()
	for i := len(bBytes); i < 32; i++ {
		buf = append(buf, 0)
	}
	temp := [][]byte{
		buf,
		bBytes,
	}
	temp2 := ConcatSlices(temp)
	ret := temp2[:32]
	return ret
}

func BigIntModMul(a *big.Int, b *big.Int, mod *big.Int) *big.Int {
	a = a.Mod(a, mod)
	b = b.Mod(b, mod)
	temp := new(big.Int)
	temp.Mul(a, b)
	temp.Mod(temp, mod)
	return temp
}

func BigIntModAdd(a *big.Int, b *big.Int, mod *big.Int) *big.Int {
	a = a.Mod(a, mod)
	b = b.Mod(b, mod)
	temp := new(big.Int)
	temp.Add(a, b)
	temp.Mod(temp, mod)
	return temp
}
