package eddsa

import (
	"bytes"
	"crypto"
	cryptorand "crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"
	"main/internal/utils"
	"math/big"
	"strconv"
)

const (
	// PublicKeySize is the size, in bytes, of public keys as used in this package.
	PublicKeySize = 32
	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
	PrivateKeySize = 64
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = 64
	// SeedSize is the size, in bytes, of private key seeds. These are the private key representations used by RFC 8032.
	SeedSize = 32

	SECURITY_BITS = 256
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey []byte

// Any methods implemented on PublicKey might need to also be implemented on
// PrivateKey, as the latter embeds the former and will expose its methods.

// Equal reports whether pub and x have the same value.
func (pub PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(PublicKey)
	if !ok {
		return false
	}
	return bytes.Equal(pub, xx)
}

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey []byte

// Public returns the PublicKey corresponding to priv.
func (priv PrivateKey) Public() crypto.PublicKey {
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, priv[32:])
	return PublicKey(publicKey)
}

// Equal reports whether priv and x have the same value.
func (priv PrivateKey) Equal(x crypto.PrivateKey) bool {
	xx, ok := x.(PrivateKey)
	if !ok {
		return false
	}
	return bytes.Equal(priv, xx)
}

// Seed returns the private key seed corresponding to priv. It is provided for
// interoperability with RFC 8032. RFC 8032's private keys correspond to seeds
// in this package.
func (priv PrivateKey) Seed() []byte {
	seed := make([]byte, SeedSize)
	copy(seed, priv[:32])
	return seed
}

// GenerateKey generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rand io.Reader) (*Keypair, error) {
	if rand == nil {
		rand = cryptorand.Reader
	}

	seed := make([]byte, SeedSize)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, err
	}

	keypair := NewKeyFromSeed(seed)
	return &keypair, nil
}

// NewKeyFromSeed calculates a private key from a seed. It will panic if
// len(seed) is not SeedSize. This function is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this
// package.
func NewKeyFromSeed(seed []byte) Keypair {
	// Outline the function body so that the returned key can be stack-allocated.
	privateKey := make([]byte, PrivateKeySize)
	keypair := newKeyFromSeed(privateKey, seed)
	return keypair
}

func newKeyFromSeed(privateKey, seed []byte) Keypair {
	if l := len(seed); l != SeedSize {
		panic("eddsa: bad seed length: " + strconv.Itoa(l))
	}
	ecPoint := ECPointGenerator()
	ecPointBytes := [32]byte{}
	ecPoint.Ge.ToBytes(&ecPointBytes)
	fmt.Println("ecPoint=", ecPointBytes)
	h := sha512.Sum512(seed)
	fmt.Println("sha512 hash=", h)

	prefix := make([]byte, 32)
	privateKey = make([]byte, 32)

	copy(prefix, h[32:])

	copy(privateKey, h[0:32])
	privateKey[0] &= 248
	privateKey[31] &= 63
	privateKey[31] |= 64

	privateKeyBN := new(big.Int).SetBytes(privateKey)
	privateKeyScalar := ECSFromBigInt(privateKeyBN)
	fmt.Print("private key 2=")
	privateKeyScalar.Print()

	prefixBN := new(big.Int).SetBytes(prefix)
	prefixScalar := ECSFromBigInt(prefixBN)
	publicKey := ecPoint.ECPMul(&privateKeyScalar.Fe)

	publicKeyBytes := [32]byte{}
	publicKey.Ge.ToBytes(&publicKeyBytes)
	fmt.Println("publicKey=", publicKeyBytes)

	return Keypair{
		PublicKey: *publicKey,
		ExtendedPrivateKey: ExpendedPrivateKey{
			Prefix:     prefixScalar,
			PrivateKey: privateKeyScalar,
		},
	}
}

func KeyAggregationN(pks *[]Ed25519Point, partyIdx uint8) *KeyAgg {
	bn1 := new(big.Int).SetUint64(1)
	xCoorVec := []big.Int{}
	for i, pk := range *pks {
		temp := pk.BytesCompressedToBigInt()
		fmt.Println("pks[] i=", i, ".bytes_compressed_to_big_int=", temp.String())
		xCoorVec = append(xCoorVec, *temp)
	}
	fmt.Println("x_coor_vec=", utils.BigIntSliceToString(&xCoorVec))

	hashVec := []big.Int{}
	for _, pk := range xCoorVec {
		vec := []big.Int{}
		vec = append(vec, *bn1)
		vec = append(vec, pk)
		for i, _ := range *pks {
			vec = append(vec, xCoorVec[i])
		}
		fmt.Println("hash_vec xx=", utils.BigIntSliceToString(&vec))
		// put all bytes together
		bytes := []byte{}
		for _, v := range vec {
			for _, b := range v.Bytes() {
				bytes = append(bytes, b)
			}
		}
		h := sha512.Sum512(bytes)
		hashVec = append(hashVec, *new(big.Int).SetBytes(h[:]))
	}
	fmt.Println("hash_vec=", utils.BigIntSliceToString(&hashVec))

	apkVec := []*Ed25519Point{}
	for i := 0; i < len(*pks); i++ {
		pk := (*pks)[i]
		hash := hashVec[i]
		hashT := ECSFromBigInt(&hash)
		a_i := pk.ECPMul(&hashT.Fe)
		apkVec = append(apkVec, a_i)
		fmt.Println("xx hash=", hash.String(), " hash_t=", hashT.ToString(), " a_i=", a_i.ToString())
	}

	sum := apkVec[0]
	for i := 1; i < len(apkVec); i++ {
		sum = sum.ECPAddPoint(&apkVec[i].Ge)
	}
	hash := ECSFromBigInt(&hashVec[partyIdx])
	keyAgg := KeyAgg{
		Apk:  *sum,
		Hash: hash,
	}
	return &keyAgg
}
