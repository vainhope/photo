package config

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)

func InitDb() *gorm.DB {
	var database = Setting.Database

	dataBaseUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", database.User, database.Password, database.Host, database.Schemas)
	database.Password = ""
	data, _ := json.Marshal(database)
	log.Printf("web config is %s", string(data))
	db, err := gorm.Open("mysql", dataBaseUrl)
	if err != nil {
		panic("connect to database fail")
	}
	//数据库别名
	db.SingularTable(true)
	//开启sql日志
	db.LogMode(true)
	return db
}

//注入db到上下文中
func InjectDB(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("db", db)
		context.Next()
	}
}
