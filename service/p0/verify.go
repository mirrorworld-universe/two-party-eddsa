package p0

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"main/internal/eddsa"
	"main/utils"
	"math/big"
)

func Verify(msg *string, R *string, s *string, publicKey *string) bool {
	msgByte := utils.StringToBytes(msg)
	msgHash := sha256.Sum256(msgByte)

	RDecoded, err := hex.DecodeString(*R)
	if err != nil {
		panic(errors.New("cannot decode R"))
	}
	sDecoded, err := hex.DecodeString(*s)
	if err != nil {
		panic(errors.New("cannot decode s"))
	}
	sDecoded = utils.ReverseByteSlice(sDecoded)

	publicKeyDecoded, err := hex.DecodeString(*publicKey)
	if err != nil {
		panic(errors.New("cannot decode publicKey"))
	}

	eight := eddsa.ECSFromBigInt(new(big.Int).SetInt64(8))
	eightInverse := eight.ModInvert()

	RDecoded32 := [32]byte{}
	for i, v := range RDecoded {
		RDecoded32[i] = v
	}
	r := eddsa.ECPFromBytes(&RDecoded32)
	r = r.ECPMul(&eightInverse.Fe)
	sBN := new(big.Int).SetBytes(sDecoded[0:32])
	sFe := eddsa.ECSFromBigInt(sBN)

	pubkey := eddsa.ECPFromBytes((*[32]byte)(publicKeyDecoded))
	pubkey = pubkey.ECPMul(&eightInverse.Fe)

	sig := eddsa.Signature{
		R:      *r,
		SmallS: sFe,
	}

	msgHash2 := msgHash[:]
	return eddsa.Verify(&sig, &msgHash2, pubkey)
}
