package db

import (
	"gorm.io/gorm"
	"main/global"
	"main/middleware/dao"
)

type MPCWallet struct {
	gorm.Model

	UserId       string `json:"user_id" gorm:"not null;index:user_id__party_id,priority:1"`
	PartyIdx     int8   `json:"party_idx" gorm:"not null;index:user_id__party_id,priority:2"`
	SeedBN       string `json:"seed_bn" gorm:"type:varchar(1000);not null"`
	KeyAggAPKBN  string `json:"key_agg_apk_bn" gorm:"type:varchar(1000);not null"`
	KeyAggHashBN string `json:"key_agg_hash_bn" gorm:"type:varchar(1000);not null"`
}

func (m MPCWallet) TableName() string {
	return "mpc_wallet"
}

func (f *MPCWallet) Create() error {
	db := dao.GetDbEngine()
	err := db.Model(&MPCWallet{}).Create(&f).Error
	if err != nil {
		global.LogDbMessage(f.TableName(), global.DbActionCreate, err.Error())
		return err
	}
	return nil
}

func (f *MPCWallet) Update() error {
	db := dao.GetDbEngine()
	err := db.Save(&f).Error
	if err != nil {
		global.LogDbMessage(f.TableName(), global.DbActionUpdate, err.Error())
		return err
	}
	return nil
}
