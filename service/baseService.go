package service

import (
	"github.com/jinzhu/gorm"
	"log"
)

func BaseAdd(db *gorm.DB, account interface{}) error {
	tx := db.Begin()
	if err := tx.Create(account).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	} else {
		tx.Commit()
	}
	return nil
}
