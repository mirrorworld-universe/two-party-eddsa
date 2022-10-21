package p0

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/levigross/grequests"
	"io/ioutil"
	"main/global"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/eddsa"
	"main/model/rest"
	"main/utils"
	"math/big"
	"time"
)

func SignRound1(userId *string, msg *string, clientKeypair *eddsa.Keypair, keyAgg *eddsa.KeyAgg) (*string, *string, *error) {
	// round 1
	//datat, err := base64.StdEncoding.DecodeString(*msg)

	// old way to do
	msgHash := sha256.Sum256([]byte(*msg))

	// msgHash from bigint
	//msgBN, _ := new(big.Int).SetString(*msg, 10)
	//msgHash := sha256.Sum256(msgBN.Bytes())
	//msgHash := msgBN.Bytes()

	println("msgbytes=", utils.BytesToStr(msgHash[:]))
	//msgHash := datat
	println("msgHash=", new(big.Int).SetBytes(msgHash[:]).String())

	clientEphemeralKey, clientSignFirstMsg, clientSignSecondMsg := eddsa.CreateEphemeralKeyAndCommit(clientKeypair, msgHash[:])
	println("clientEphemeralKey=", clientEphemeralKey.ToString(), ", clientSignFirstMsg=", clientSignFirstMsg.ToString()+", clientSignSecondMsg=", clientSignSecondMsg.ToString())

	// send request to P1 to get commitment
	url := global.Config.Base.P1Url + "/p1/sign_round1"
	data := map[string]interface{}{
		"user_id":          userId,
		"client_pubkey_bn": clientKeypair.PublicKey.BytesCompressedToBigInt().String(),
		"msg_hash_bn":      new(big.Int).SetBytes(msgHash[:]).String(),
	}

	var resp rest.P1SignRound1Response
	response, err := grequests.Post(url, &grequests.RequestOptions{
		JSON:           data,
		RequestTimeout: time.Second * 5,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	body, err := ioutil.ReadAll(response.RawResponse.Body)
	defer response.RawResponse.Body.Close()

	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(errors.New("error parse p1Round1 response"))
	}
	serverSignFirstMsgCommitment, _ := new(big.Int).SetString(resp.ServerSignFirstMsgCommitmentBN, 10)
	println("[P0SignRound1] p1_sign_round1 resp, ServerSignFirstMsgCommitmentBN=", resp.ServerSignFirstMsgCommitmentBN)

	// p1 round2
	url = global.Config.Base.P1Url + "/p1/sign_round2"
	data = map[string]interface{}{
		"user_id":                             userId,
		"client_pubkey_bn":                    clientKeypair.PublicKey.BytesCompressedToBigInt().String(),
		"msg_hash_bn":                         new(big.Int).SetBytes(msgHash[:]).String(),
		"client_sign_first_msg_commitment_bn": clientSignFirstMsg.Commitment.String(),
		"client_sign_second_msg_r_bn":         clientSignSecondMsg.R.BytesCompressedToBigInt().String(),
		"client_sign_second_msg_bf_32_bn":     new(big.Int).SetBytes(utils.BigintToBytes32(&clientSignSecondMsg.BlindFactor)).String(),
	}
	var resp2 rest.P1SignRound2Response
	response, err = grequests.Post(url, &grequests.RequestOptions{
		JSON:           data,
		RequestTimeout: time.Second * 5,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	body, err = ioutil.ReadAll(response.RawResponse.Body)
	defer response.RawResponse.Body.Close()

	err = json.Unmarshal(body, &resp2)
	if err != nil {
		panic(errors.New("error parse p1Round1 response"))
	}
	println("[P0SignRound1]p1_sign_round2 resp, ServerSignSecondMsgR=", resp2.ServerSignSecondMsgR, ", ServerSignSecondMsgBF32=", resp2.ServerSignSecondMsgBF32, " ServerSigRBN=", resp2.ServerSigRBN, " ServerSigSmallSBN=", resp2.ServerSigSmallSBN)

	serverSignSecondMsgRBN, _ := new(big.Int).SetString(resp2.ServerSignSecondMsgR, 10)
	serverSignSecondMsgR := eddsa.NewECPSetFromBN(serverSignSecondMsgRBN)

	ServerSignSecondMsgBF32BN, _ := new(big.Int).SetString(resp2.ServerSignSecondMsgBF32, 10)

	serverSigRBN, _ := new(big.Int).SetString(resp2.ServerSigRBN, 10)
	serverSigSmallSBN, _ := new(big.Int).SetString(resp2.ServerSigSmallSBN, 10)
	serverSig := eddsa.Signature{
		R:      *eddsa.NewECPSetFromBN(serverSigRBN),
		SmallS: *eddsa.NewECSSetFromBN(serverSigSmallSBN),
	}
	println("[P0SignRound1]serverSig=", serverSig.ToString())
	// check commiment
	isCommMatch := eddsa.CheckCommitment(
		serverSignSecondMsgR,
		ServerSignSecondMsgBF32BN,
		serverSignFirstMsgCommitment,
	)

	if !isCommMatch {
		panic(errors.New("commitment not match"))
	}

	// round 3
	ri := []eddsa.Ed25519Point{
		*serverSignSecondMsgR,
		clientSignSecondMsg.R,
	}
	rTot := eddsa.SigGetRTot(ri)
	println("[P0SignRound1] rTot=", rTot.ToString())

	msgHash2 := msgHash[:]
	k := eddsa.SigK(rTot, &keyAgg.Apk, &msgHash2)
	println("[P0SignRound1] k=", k.ToString(), " keyAgg=", keyAgg.ToString(), " msgHash=", new(big.Int).SetBytes(msgHash2).String())

	s2 := eddsa.PartialSign(
		&clientEphemeralKey.SmallR,
		clientKeypair,
		&k,
		&keyAgg.Hash,
		rTot,
	)
	println("[P0SignRound1] s2=", s2.ToString())

	s := []eddsa.Signature{
		serverSig,
		s2,
	}
	sig := eddsa.AddSignatureParts(s)
	RBytes := [32]byte{}
	sig.R.Ge.ToBytes(&RBytes)
	sBytes := [32]byte{}
	edwards25519.FeToBytes(&sBytes, &sig.SmallS.Fe)
	println("[P0SignRound1] sig=", sig.ToString(), " R: ", hex.EncodeToString(RBytes[:]), " s:", hex.EncodeToString(sBytes[:]))

	// final verify
	R := hex.EncodeToString(RBytes[:])
	smallS := hex.EncodeToString(sBytes[:])
	if isMatch := eddsa.Verify(&sig, &msgHash2, &keyAgg.Apk); isMatch {
		return &R, &smallS, nil
	}

	err = errors.New("invalid_signature")
	return nil, nil, &err
}

func Sign(msg *string, clientKeypair *eddsa.Keypair, keyAgg *eddsa.KeyAgg) {

	//clientKeypair, keyAgg := tempLoadKey()
	println("clientKeyPair=", clientKeypair.ToString(), " keyagg=", keyAgg.ToString())

	// round 1
	msgHash := sha256.Sum256([]byte(*msg))
	println("msgHash=", new(big.Int).SetBytes(msgHash[:]).String())

	clientEphemeralKey, clientSignFirstMsg, clientSignSecondMsg := eddsa.CreateEphemeralKeyAndCommit(clientKeypair, msgHash[:])
	println("clientEphemeralKey=", clientEphemeralKey.ToString(), ", clientSignFirstMsg=", clientSignFirstMsg.ToString()+", clientSignSecondMsg=", clientSignSecondMsg.ToString())

	// now send clientSignFirstMsg, msgHash, client public key to p1, and receive serverFirstSignMsg
	serverCommitment, _ := new(big.Int).SetString("84931746524459149992060349634228453990530694124359495037784716096273864068584", 10)
	serverSignFirstMsg := eddsa.SignFirstMsg{
		Commitment: *serverCommitment,
	}
	println("serverSignFirstMsg=", serverSignFirstMsg.ToString())

	// round 2
	// send clientSecondSignMsg to p1, get serverSignSecondMsg{R, blindFactor}
	eight := eddsa.ECSFromBigInt(new(big.Int).SetInt64(8))
	eightInverse := eight.ModInvert()
	serverSignSecondMsgRBytes := [32]byte{
		142, 144, 114, 134, 190, 107, 127, 90,
		212, 252, 156, 101, 121, 82, 106, 155,
		187, 60, 75, 220, 240, 209, 132, 217,
		100, 78, 252, 14, 20, 73, 153, 54,
	}
	serverSignSecondMsgR := eddsa.ECPFromBytes(&serverSignSecondMsgRBytes)
	serverSignSecondMsgR = serverSignSecondMsgR.ECPMul(&eightInverse.Fe)

	temp1 := [32]byte{
		169, 43, 89, 150, 255, 113, 182, 143,
		232, 177, 192, 27, 76, 61, 36, 72,
		121, 68, 213, 61, 241, 206, 20, 165,
		112, 33, 80, 6, 72, 206, 30, 83,
	}
	serverSignSecondMsgBF := new(big.Int).SetBytes(temp1[:])
	println("round2, server_sign_second_msg_R=", serverSignSecondMsgR.ToString(), ", serverSignSecondMsgBF=", serverSignSecondMsgBF.String())

	temp2 := [32]byte{
		29, 217, 173, 145, 214, 96, 225, 4,
		252, 2, 4, 62, 123, 190, 12, 48,
		63, 27, 252, 28, 1, 38, 137, 171,
		140, 45, 56, 196, 174, 107, 224, 231,
	}
	serverSigR := eddsa.ECPFromBytes(&temp2)
	serverSigR = serverSigR.ECPMul(&eightInverse.Fe)
	println("round2, serverSigR=", serverSigR.ToString())

	temp3 := [32]byte{
		124, 155, 253, 249, 189, 116, 9, 104,
		139, 154, 108, 227, 90, 150, 239, 201,
		172, 186, 250, 211, 86, 58, 200, 208,
		138, 102, 125, 137, 46, 247, 205, 10,
	}
	temp33 := utils.ReverseByteSlice(temp3[:])
	serverSigS := eddsa.ECSFromBigInt(new(big.Int).SetBytes(temp33))

	serverSignSecondMsg := eddsa.SignSecondMsg{
		R:           *serverSignSecondMsgR,
		BlindFactor: *serverSignSecondMsgBF,
	}
	serverSig := eddsa.Signature{
		R:      *serverSigR,
		SmallS: serverSigS,
	}
	println("round2, serverSignSecondMsg=", serverSignSecondMsg.ToString(), " serverSig=", serverSig.ToString())

	// check commiment
	isCommMatch := eddsa.CheckCommitment(
		&serverSignSecondMsg.R,
		&serverSignSecondMsg.BlindFactor,
		&serverSignFirstMsg.Commitment,
	)

	if !isCommMatch {
		panic(errors.New("commitment not match"))
	}

	// round 3
	ri := []eddsa.Ed25519Point{
		*serverSignSecondMsgR,
		clientSignSecondMsg.R,
	}
	rTot := eddsa.SigGetRTot(ri)
	println("rTot=", rTot.ToString())

	msgHash2 := msgHash[:]
	k := eddsa.SigK(rTot, &keyAgg.Apk, &msgHash2)
	println("k=", k.ToString())

	s2 := eddsa.PartialSign(
		&clientEphemeralKey.SmallR,
		clientKeypair,
		&k,
		&keyAgg.Hash,
		rTot,
	)
	println("s2=", s2.ToString())

	s := []eddsa.Signature{
		serverSig,
		s2,
	}
	sig := eddsa.AddSignatureParts(s)
	RBytes := [32]byte{}
	sig.R.Ge.ToBytes(&RBytes)
	sBytes := [32]byte{}
	edwards25519.FeToBytes(&sBytes, &sig.SmallS.Fe)
	println("sig=", sig.ToString(), " R: ", hex.EncodeToString(RBytes[:]), " s:", hex.EncodeToString(sBytes[:]))

	// final verify
	eddsa.Verify(&sig, &msgHash2, &keyAgg.Apk)

	//clientPublicKeyBytes := [32]byte{}
	//clientKeypair.PublicKey.Ge.ToBytes(&clientPublicKeyBytes)
	//
	//temp := [][]byte{
	//	[]byte{2},
	//	utils.BigintToBytes32(&clientSignFirstMsg.Commitment),
	//	msgHash[:32],
	//	clientPublicKeyBytes[:],
	//}
	//
	//// send to p1
	//
	//buf := make([]byte, 32)
	//serverCommitmentBytes := []byte{}
	//serverCommitment := new(big.Int).SetBytes(serverCommitmentBytes)
	//serverSignFirstMsg := eddsa.SignFirstMsg{
	//	Commitment: *serverCommitment,
	//}
	//
	//// round 2
	//clientSignSecondMsgBytes := [32]byte{}
	//clientSignSecondMsg.R.Ge.ToBytes(&clientSignSecondMsgBytes)
	//temp = [][]byte{
	//	clientSignSecondMsgBytes[:],
	//	utils.BigintToBytes32(&clientSignSecondMsg.BlindFactor),
	//}
	//buf = utils.ConcatSlices(temp)
	// send to p1

}
