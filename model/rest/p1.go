package rest

type P1KeygenRound1Req struct {
	ClientPubkeyBN string `json:"client_pubkey_bn" vd:"len($)>0"`
	ServerSKSeed   string `json:"server_sk_seed"`
}

type P1KeygenRound1Response struct {
	ServerPubkeyBN string `json:"server_pubkey_bn" vd:"len($)>0"`
}

type P1SignRound1Req struct {
	ClientSignFirstMsgCommitmentBN string `json:"client_sign_first_msg_commitment_bn" vd:"len($)>0"`
	MsgHashBN                      string `json:"msg_hash_bn" vd:"len($)>0"`
	ClientPubkeyBN                 string `json:"client_pubkey_bn" vd:"len($)>0"`
}

type P1SignRound1Response struct {
	ServerSignFirstMsgCommitmentBN string `json:"server_sign_first_msg_commitment_bn" vd:"len($)>0"`
}
