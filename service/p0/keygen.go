package p0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"io/ioutil"
	"main/global"
	"main/internal/agl_ed25519/edwards25519"
	"main/internal/eddsa"
	"main/model/rest"
	"math/big"
	"time"
)

func KeyGenRound1NoSeed() (*eddsa.Keypair, *eddsa.KeyAgg) {
	ecsRndBytes := [32]byte{}
	edwards25519.FeToBytes(&ecsRndBytes, &eddsa.ECSNewRandom().Fe)
	sKSeed := new(big.Int).SetBytes(ecsRndBytes[:])

	return keyGenRound1Internal(sKSeed, nil)
}

func KeyGenRound1FromSeed(clientSKSeed *big.Int) (*eddsa.Keypair, *eddsa.KeyAgg) {
	return keyGenRound1Internal(clientSKSeed, nil)
}

func KeyGenRound1FromBothSeed(clientSKSeed *big.Int, serverSKSeed *big.Int) (*eddsa.Keypair, *eddsa.KeyAgg) {
	return keyGenRound1Internal(clientSKSeed, serverSKSeed)
}

func keyGenRound1Internal(clientSKSeed *big.Int, serverSKSeed *big.Int) (*eddsa.Keypair, *eddsa.KeyAgg) {
	// generate client keypair
	clientKeypair := eddsa.CreateKeyPairFromSeed(clientSKSeed)
	clientPublicKeyBytes := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPublicKeyBytes)

	// ask for server public key
	data := map[string]interface{}{
		"client_pubkey_bn": new(big.Int).SetBytes(clientPublicKeyBytes[:]).String(),
	}
	if serverSKSeed != nil {
		data["server_sk_seed"] = serverSKSeed.String()
	}
	response, err := grequests.Post("http://localhost:3000/p1/keygen_round1", &grequests.RequestOptions{
		JSON:           data,
		RequestTimeout: time.Second * 5,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	var resp rest.P1KeygenRound1Response
	body, err := ioutil.ReadAll(response.RawResponse.Body)
	defer response.RawResponse.Body.Close()

	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(errors.New("error parse p1Round1 response"))
	}

	serverPubkeyBN, _ := new(big.Int).SetString(resp.ServerPubkeyBN, 10)
	serverPubkey := eddsa.ECPFromBytes((*[32]byte)(serverPubkeyBN.Bytes()))

	// start aggregate
	pks := []eddsa.Ed25519Point{
		*serverPubkey,           // partyIdx=0
		clientKeypair.PublicKey, // partyIdx=1
	}
	keyAgg := eddsa.KeyAggregationN(&pks, global.PARTY_INDEX_P0)
	aggPubKeyBytes := [32]byte{}
	keyAgg.Apk.Ge.ToBytes(&aggPubKeyBytes)
	//fmt.Println("aggregated_pukey=", hex.EncodeToString(aggPubKeyBytes[:]))
	//fmt.Println("key_agg=", keyAgg.ToString())
	return clientKeypair, keyAgg
}

func KeyGen() (*eddsa.Keypair, *eddsa.KeyAgg) {
	fmt.Println("*************Client*************")
	rnd, _ := new(big.Int).SetString("5266194697103632731894445446481908111422432681065623019013231350200571873746", 10)
	clientKeypair := eddsa.CreateKeyPairFromSeed(rnd)
	clientPublicKeyBytes := [32]byte{}
	clientKeypair.PublicKey.Ge.ToBytes(&clientPublicKeyBytes)

	fmt.Println("*************Server*************")
	rnd, _ = new(big.Int).SetString("1276567075174267627823301091809777026200725024551313144625936661005557002592", 10)
	serverKeypair := eddsa.CreateKeyPairFromSeed(rnd)
	serverPublicKeyBytes := [32]byte{}
	serverKeypair.PublicKey.Ge.ToBytes(&serverPublicKeyBytes)

	// start aggregate
	pks := []eddsa.Ed25519Point{
		serverKeypair.PublicKey, // partyIdx=0
		clientKeypair.PublicKey, // partyIdx=1
	}
	keyAgg := eddsa.KeyAggregationN(&pks, global.PARTY_INDEX_P0)
	aggPubKeyBytes := [32]byte{}
	keyAgg.Apk.Ge.ToBytes(&aggPubKeyBytes)
	fmt.Println("aggregated_pukey=", hex.EncodeToString(aggPubKeyBytes[:]))
	fmt.Println("key_agg=", keyAgg.ToString())

	// @TODO save to db
	return clientKeypair, keyAgg

	// continue sign process
}
