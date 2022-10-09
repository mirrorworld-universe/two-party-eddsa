package rest

type P0KeygenRound1Response struct {
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required"`
	KeyAgg         string `json:"key_agg" binding:"required"`
}
