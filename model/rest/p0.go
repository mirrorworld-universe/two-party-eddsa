package rest

// key gen
type P0KeygenRound1Req struct {
	ClientSKSeed string `json:"client_sk_seed"`
	ServerSKSeed string `json:"server_sk_seed"` // test purpose. Not in prd env
}

type P0KeygenRound1Response struct {
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required,validbn"`
	KeyAgg         string `json:"key_agg" binding:"required"`
	KeyAggBs58     string `json:"key_agg_bs58" binding:"required"`
	UserId         string `json:"user_id" binding:"required"`
}

// get key gen by user id
type P0GetKeyGenByUserIdResponse struct {
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required,validbn"`
	KeyAgg         string `json:"key_agg" binding:"required"`
	KeyAggBs58     string `json:"key_agg_bs58" binding:"required"`
}

// sign
type P0SignRound1Req struct {
	UserId string `json:"user_id" binding:"required"`
	Msg    string `json:"msg" binding:"required"`
}

type P0SignRound1Response struct {
	R      string `json:"R" binding:"required"`
	SmallS string `json:"s" binding:"required"`
}

// verify
type P0VerifyReq struct {
	Msg    string `json:"msg" binding:"required"`
	R      string `json:"R" binding:"required"`
	SmallS string `json:"s" binding:"required"`
	KeyAgg string `json:"key_agg" binding:"required"`
}

type P0VerifyResponse struct {
	IsValid bool `json:"is_valid" binding:"required"`
}
