package p0

import (
	"crypto/sha512"
	"main/internal/eddsa"
)

func Sign(msg *string, clientKeypair *eddsa.Keypair, keyAgg *eddsa.KeyAgg) {
	// round 1
	msgHash := sha512.Sum512([]byte(*msg))

}
