package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"main/model/rest"
	"net/http"
	"net/http/httptest"
)

func (t *SuiteTest) TestPing() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/p0/test", nil)
	t.router.ServeHTTP(w, req)

	//assertion
	assert.Equal(t.T(), http.StatusOK, w.Code)
}

func verifyKeyGen(t *SuiteTest, v *TestInput) *string {
	w := httptest.NewRecorder()
	requestBody := map[string]string{
		"client_sk_seed": v.ClientSKSeed,
		"server_sk_seed": v.ServerSKSeed,
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/p0/keygen_round1", bytes.NewBuffer(jsonData))
	t.router.ServeHTTP(w, req)

	// parse result
	var resp rest.P0KeygenRound1Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusOK, w.Code)
	assert.Equal(t.T(), v.ExpectedClientPubkeyBN, resp.ClientPubkeyBN)
	assert.Equal(t.T(), v.ExpectedKeyAgg, resp.KeyAgg)
	assert.True(t.T(), len(resp.UserId) > 0)

	return &resp.UserId
}

func verifySignMsg(t *SuiteTest, userId *string, v *TestInput) {
	w := httptest.NewRecorder()
	requestBody := map[string]string{
		"user_id": *userId,
		"msg":     v.Msg,
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/p0/sign_round1", bytes.NewBuffer(jsonData))
	t.router.ServeHTTP(w, req)

	// parse result
	var resp rest.P0SignRound1Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusOK, w.Code)
	assert.Equal(t.T(), v.ExpectedR, resp.R)
	assert.Equal(t.T(), v.ExpectedSmallS, resp.SmallS)
}

func (t *SuiteTest) TestP0KeyGen() {

	v := TestInput{
		ClientSKSeed:           "5266194697103632731894445446481908111422432681065623019013231350200571873746",
		ServerSKSeed:           "1276567075174267627823301091809777026200725024551313144625936661005557002592",
		ExpectedClientPubkeyBN: "25044372729105238728613876994056928120890707228207216431325866756268369637254",
		ExpectedKeyAgg:         "790c23f4a2f065fa4cebf77a005f75ad7a528c8de4ca64e4e5c681c17663514e",
	}
	verifyKeyGen(t, &v)
}

func (t *SuiteTest) TestP0KeyGenAndSign() {
	for _, v := range testcases {

		// step 1. key gen
		userId := verifyKeyGen(t, &v)

		// step 2: sign msg
		verifySignMsg(t, userId, &v)
	}
}
