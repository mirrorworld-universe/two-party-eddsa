package rest

// P1 KeyGen
type P1KeygenRound1Req struct {
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required,validbn"`
	ServerSKSeed   string `json:"server_sk_seed"`
}

type P1KeygenRound1Response struct {
	ServerPubkeyBN string `json:"server_pubkey_bn" binding:"required,validbn"`
}

// P1 Sign round1
type P1SignRound1Req struct {
	MsgHashBN      string `json:"msg_hash_bn" binding:"required,validbn"`
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required,validbn"`
}

type P1SignRound1Response struct {
	ServerSignFirstMsgCommitmentBN string `json:"server_sign_first_msg_commitment_bn" binding:"required,validbn"`
}

// P1 Sign round2
type P1SignRound2Req struct {
	ClientSignFirstMsgCommitmentBN string `json:"client_sign_first_msg_commitment_bn" binding:"required,validbn"`
	ClientSignSecondMsgRBN         string `json:"client_sign_second_msg_r_bn" binding:"required,validbn"`
	ClientSignSecondMsgBFBN        string `json:"client_sign_second_msg_bf_bn" binding:"required,validbn"`
}

type P1SignRound2Response struct {
	ServerSignFirstMsgCommitmentBN string `json:"server_sign_first_msg_commitment_bn" binding:"required,validbn"`
}
