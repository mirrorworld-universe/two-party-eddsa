package rest

// P1 KeyGen
type P1KeygenRound1Req struct {
	UserId         string `json:"user_id" binding:"required"`
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required,validbn"`
	ServerSKSeed   string `json:"server_sk_seed"`
}

type P1KeygenRound1Response struct {
	ServerPubkeyBN string `json:"server_pubkey_bn" binding:"required,validbn"`
}

// P1 Sign round1
type P1SignRound1Req struct {
	UserId         string `json:"user_id" binding:"required"`
	MsgHashBN      string `json:"msg_hash_bn" binding:"required,validbn"`
	ClientPubkeyBN string `json:"client_pubkey_bn" binding:"required,validbn"`
}

type P1SignRound1Response struct {
	ServerSignFirstMsgCommitmentBN string `json:"server_sign_first_msg_commitment_bn" binding:"required,validbn"`
}

// P1 Sign round2
type P1SignRound2Req struct {
	UserId                         string `json:"user_id" binding:"required"`
	ClientPubkeyBN                 string `json:"client_pubkey_bn" binding:"required,validbn"`
	MsgHashBN                      string `json:"msg_hash_bn" binding:"required,validbn"`
	ClientSignFirstMsgCommitmentBN string `json:"client_sign_first_msg_commitment_bn" binding:"required,validbn"`
	ClientSignSecondMsgRBN         string `json:"client_sign_second_msg_r_bn" binding:"required,validbn"`
	ClientSignSecondMsgBF32BN      string `json:"client_sign_second_msg_bf_32_bn" binding:"required,validbn"`
}

type P1SignRound2Response struct {
	ServerSignSecondMsgR    string `json:"server_sign_second_msg_r" binding:"required"`
	ServerSignSecondMsgBF32 string `json:"server_sign_second_msg_bf_32" binding:"required"`
	ServerSigRBN            string `json:"server_sig_r_bn" binding:"required,validbn"`
	ServerSigSmallSBN       string `json:"server_sig_small_s_bn" binding:"required,validbn"`
}
