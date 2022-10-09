package request

type P1KeygenRound1Req struct {
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required"`
	ServerSKSeed   string `json:"server_sk_seed"`
}

type P1KeygenRound1Response struct {
	ServerPubkeyBN string `json:"server_pubkey_bn" binding:"required"`
}
