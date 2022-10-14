package tests

type TestInput struct {
	ClientSKSeed           string `json:"client_sk_seed"`
	ServerSKSeed           string `json:"server_sk_seed"` // test purpose. Not in prd env
	ExpectedClientPubkeyBN string `json:"expected_client_pubkey_bn" binding:"validbn"`
	ExpectedKeyAgg         string `json:"expected_key_agg"`
	Msg                    string `json:"msg"`
	ExpectedR              string `json:"R"`
	ExpectedSmallS         string `json:"s"`
}

var testcases = []TestInput{
	TestInput{
		ClientSKSeed:           "5266194697103632731894445446481908111422432681065623019013231350200571873746",
		ServerSKSeed:           "1276567075174267627823301091809777026200725024551313144625936661005557002592",
		ExpectedClientPubkeyBN: "25044372729105238728613876994056928120890707228207216431325866756268369637254",
		ExpectedKeyAgg:         "790c23f4a2f065fa4cebf77a005f75ad7a528c8de4ca64e4e5c681c17663514e",
		Msg:                    "hello",
		ExpectedR:              "1dd9ad91d660e104fc02043e7bbe0c303f1bfc1c012689ab8c2d38c4ae6be0e7",
		ExpectedSmallS:         "7bf0d2eb8027a65988c43a4c79e70f3ab67eadf1a8a852b5cf34ef1ace192407",
	},
}
