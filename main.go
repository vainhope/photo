package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"goweb/config"
	"goweb/route"
)

func main() {
	logger := config.CreateLogFile()

	db := config.InitDb()

	defer db.Close()

	r := gin.Default()

	r = config.InitServer(r)

	//必须在路由之前注入db到上下文中
	r.Use(config.InjectDB(db))

	//注入log
	r.Use(config.InjectLog(logger))

	route.Router(r)

	err := r.Run(fmt.Sprintf(":%s", config.Setting.Server.Port))

	if err != nil {
		logger.Println(err.Error())
	}
}
