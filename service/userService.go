package service

import (
	"github.com/jinzhu/gorm"
	"goweb/entity"
	"log"
)

func GreaterThan(db *gorm.DB) *gorm.DB {
	return db.Where("id > ?", 100)
}

func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(GreaterThan).Where("id IN (?)", status)
	}
}

func FindByName(db *gorm.DB, userName string) entity.User {
	var account entity.User
	db.Table("t_user").Where(map[string]interface{}{
		"userName": userName,
		"deleted":  0,
	}).First(&account)
	return account
}

func FindById(db *gorm.DB, id string) entity.User {
	var account entity.User
	db.Table("t_user").Scopes(OrderStatus([]string{id})).Limit(1).Scan(&account)
	return account
}

func FindInfoWithIds(db *gorm.DB, ids []string) []entity.User {
	var accounts []entity.User
	db.Raw("SELECT * FROM t_user WHERE id IN (?)", ids).Scan(&accounts)
	return accounts
}

func Add(db *gorm.DB, account *entity.User) error {
	omitColumns := make([]string, 2)
	//忽略字段
	if "" == account.Email {
		omitColumns = append(omitColumns, "email")
	}
	if "" == account.Email {
		omitColumns = append(omitColumns, "phone")
	}
	tx := db.Begin().Omit(omitColumns...)
	if err := tx.Create(account).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	} else {
		tx.Commit()
	}
	return nil
}
