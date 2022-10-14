package dao

import (
	"gorm.io/gorm"
)

var dbEngine *gorm.DB

func GetDbEngine() *gorm.DB {
	return dbEngine
}

func SetDbEngine(db *gorm.DB) {
	dbEngine = db
}
