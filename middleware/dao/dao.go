package dao

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Dao struct {
	DB *gorm.DB
}

func NewDao() *Dao {
	return &Dao{
		DB: GetDbEngine(),
	}
}

func NewDaoEmpty() *Dao {
	return &Dao{}
}

func (d *Dao) Health() (string, error) {
	data := fmt.Sprintf("server is running, datetime: %v", time.Now().Format("2006-01-02 15:04:05"))
	return data, nil
}
