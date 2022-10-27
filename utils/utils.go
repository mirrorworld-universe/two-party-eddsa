package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/levigross/grequests"
	"io/ioutil"
	"main/global"
	"math/big"
	"strconv"
	"time"
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
	a2 := a.Mod(a, mod)
	b2 := b.Mod(b, mod)
	temp := new(big.Int)
	temp.Add(a2, b2)
	temp.Mod(temp, mod)
	return temp
}

func BigIntToByte32(n *big.Int) *[32]byte {
	nByte := n.Bytes()
	if len(nByte) < 32 {
		// insert 0 in the begining
		bytes := [][]byte{
			make([]byte, 32-len(nByte)),
			nByte,
		}
		nByte = ConcatSlices(bytes)
	}
	ret := nByte[0:32]
	return (*[32]byte)(ret)
}

func SendReqAndParseResp[T any](url *string, data *map[string]interface{}, respSchema *T) error {
	response, err := grequests.Post(global.P1Url()+"/p1/keygen_round1", &grequests.RequestOptions{
		JSON:           data,
		RequestTimeout: time.Second * 5,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	body, err := ioutil.ReadAll(response.RawResponse.Body)
	defer response.RawResponse.Body.Close()

	err = json.Unmarshal(body, &respSchema)
	if err != nil {
		return errors.New("error parse p1Round1 response")
	}
	return nil
}
