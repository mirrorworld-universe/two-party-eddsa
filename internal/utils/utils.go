package utils

import (
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
