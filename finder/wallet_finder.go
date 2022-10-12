package finder

import (
	"main/middleware/dao"
	"main/model/db"
)

func FindP0ByUserId(userId *string) *db.MPCWallet {
	var wallet db.MPCWallet
	dao.GetDbEngine().First(&wallet, "user_id=? AND party_idx=?", userId, 0)
	return &wallet
}

func FindP1ByUserId(userId *string) *db.MPCWallet {
	var wallet db.MPCWallet
	dao.GetDbEngine().First(&wallet, "user_id=? AND party_idx=?", userId, 1)
	return &wallet
}
