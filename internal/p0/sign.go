package p0

import (
	"crypto/sha512"
	"main/internal/eddsa"
	"main/internal/utils"
	"math/big"
)

func Sign(msg *string, clientKeypair *eddsa.Keypair, keyAgg *eddsa.KeyAgg) {
	// round 1
	msgHash := sha512.Sum512([]byte(*msg))
	clientEphemeralKey, clientSignFirstMsg, clientSignSecondMsg := eddsa.CreateEphemeralKeyAndCommit(clientKeypair, msgHash[:])
	clientPublicKeyBytes := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPublicKeyBytes)

	temp := [][]byte{
		[]byte{2},
		utils.BigintToBytes32(&clientSignFirstMsg.Commitment),
		msgHash[:32],
		clientPublicKeyBytes[:],
	}

	// send to server

	buf := make([]byte, 32)
	serverCommitmentBytes := []byte{}
	serverCommitment := new(big.Int).SetBytes(serverCommitmentBytes)
	serverSignFirstMsg := eddsa.SignFirstMsg{
		Commitment: *serverCommitment,
	}

	// round 2
	clientSignSecondMsgBytes := [32]byte{}
	clientSignSecondMsg.R.Ge.ToBytes(&clientSignSecondMsgBytes)
	temp = [][]byte{
		clientSignSecondMsgBytes[:],
		utils.BigintToBytes32(&clientSignSecondMsg.BlindFactor),
	}
	buf = utils.ConcatSlices(temp)
	// send to server

}
