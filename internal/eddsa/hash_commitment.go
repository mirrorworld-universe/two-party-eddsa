package eddsa

import "golang.org/x/crypto/sha3"
import (
	"main/internal/utils"
	"math/big"
)

// calculate commitment c = H(m,r) using SHA3 CRHF. r is 256bit blinding factor, m is the commited value
type HashCommitment struct {
	Commitment  big.Int
	BlindFactor big.Int
}

func CreateCommitmentWithUserDefinedRandomness(message *big.Int, blindingFactor *big.Int) *big.Int {
	bytes := [][]byte{
		message.Bytes(),
		blindingFactor.Bytes(),
	}
	bytesAll := utils.ConcatSlices(bytes)
	h := sha3.Sum256(bytesAll)
	return new(big.Int).SetBytes(h[:])
}

func CreateCommitment(message *big.Int) *HashCommitment {
	blindFactor := utils.BigintSample(SECURITY_BITS)
	com := CreateCommitmentWithUserDefinedRandomness(message, blindFactor)
	return &HashCommitment{
		Commitment:  *com,
		BlindFactor: *blindFactor,
	}
}

func CheckCommitment(smallRToTest *Ed25519Point, blindFactor *big.Int, comm *big.Int) bool {
	computedComm := CreateCommitmentWithUserDefinedRandomness(smallRToTest.BytesCompressedToBigInt(), blindFactor)
	return computedComm == comm
}
