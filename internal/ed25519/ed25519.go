package ed25519

import (
	"bytes"
	"crypto"
	cryptorand "crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"
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
func GenerateKey(rand io.Reader) (PublicKey, PrivateKey, error) {
	if rand == nil {
		rand = cryptorand.Reader
	}

	seed := make([]byte, SeedSize)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, nil, err
	}

	privateKey := NewKeyFromSeed(seed)
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, privateKey[32:])

	return publicKey, privateKey, nil
}

// NewKeyFromSeed calculates a private key from a seed. It will panic if
// len(seed) is not SeedSize. This function is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this
// package.
func NewKeyFromSeed(seed []byte) PrivateKey {
	// Outline the function body so that the returned key can be stack-allocated.
	privateKey := make([]byte, PrivateKeySize)
	newKeyFromSeed(privateKey, seed)
	return privateKey
}

func newKeyFromSeed(privateKey, seed []byte) Keypair {
	if l := len(seed); l != SeedSize {
		panic("ed25519: bad seed length: " + strconv.Itoa(l))
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
	fmt.Println("private key 2=")
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

//type GE = Ed25519Point
//type FE = Ed25519Scalar
//
//func (e *Ed25519Scalar) toBigInt() *big.Int {
//	feBytes := [32]byte{}
//	edwards25519.FeToBytes(&feBytes, &e.Fe)
//
//	// reverse fe_bytes
//	for i, j := 0, len(feBytes)-1; i < j; i, j = i+1, j-1 {
//		feBytes[i], feBytes[j] = feBytes[j], feBytes[i]
//	}
//
//	ret := new(big.Int).SetBytes(feBytes[:])
//	return ret
//}
//
//func q() *big.Int {
//	qBytesArray := [32]byte{237, 211, 245, 92, 26, 99, 18, 88, 214, 156, 247, 162, 222, 249, 222, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16}
//	lFe := new(SK)
//	edwards25519.FeFromBytes(lFe, &qBytesArray)
//	lFeScalar := Ed25519Scalar{
//		Purpose: "q",
//		Fe:      *lFe,
//	}
//	return lFeScalar.toBigInt()
//}
//
//func from(n *big.Int) *Ed25519Scalar {
//	n_bytes := n.Bytes()
//	n_bytes_64 = n_bytes[:]
//	n_bytes_r := utils.ReverseByteSlice(n_bytes)
//	out := [32]byte{}
//	edwards25519.ScReduce(&out, &n_bytes_r)
//}
//
//func (e *Ed25519Scalar) newRandom() *Ed25519Scalar {
//	// sample_below()
//	reader := cryptorand.Reader
//	rnd_bn, _ := cryptorand.Int(reader, q())
//	bn_8 := big.NewInt(8)
//	rnd_bn_mul := new(big.Int).Mul(rnd_bn, bn_8)
//	rnd_bn_mul_8 := new(big.Int).Mod(rnd_bn_mul, q())
//	return rnd_bn_mul_mod
//}
