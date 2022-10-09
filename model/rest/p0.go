package rest

type P0KeygenRound1Req struct {
	ClientSKSeed string `json:"client_sk_seed"`
	ServerSKSeed string `json:"server_sk_seed"` // test purpose. Not in prd env
}

type P0KeygenRound1Response struct {
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required"`
	KeyAgg         string `json:"key_agg" binding:"required"`
}
