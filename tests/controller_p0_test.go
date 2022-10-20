package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"main/model/rest"
	"net/http"
	"net/http/httptest"
)

/**
Sample test pong endpoint
*/
func (t *SuiteTest) TestPing() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/p0/test", nil)
	t.router.ServeHTTP(w, req)

	//assertion
	assert.Equal(t.T(), http.StatusOK, w.Code)
}

/**
Helper method for keygen
*/
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

/**
Helper method for sign msg
*/
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
		ClientSKSeed:           "6152830912195565519395098481732012873411846110202199315509340643094565512733",
		ServerSKSeed:           "2907473004572530518387226344766367270548447211691788173429811968516240988842",
		ExpectedClientPubkeyBN: "62564454420585069468955050107550805769837960470189050038393387670380364818689",
		ExpectedKeyAgg:         "3d1f64651ed8fc0f9ff867b1abf4100945f0b774409422576276016cc2b9e131",
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
