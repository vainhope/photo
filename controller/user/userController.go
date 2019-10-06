package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"goweb/base"
	"goweb/entity"
	"goweb/service"
	"goweb/util"
	"goweb/utils"
	"strings"
)

func Login(context *gin.Context) {
	var account entity.User
	if err := context.ShouldBindJSON(&account); err != nil {
		base.Response(context, base.NOT_LOGIN, nil)
		return
	}
	if "" == account.UserName || "" == account.PassWord {
		base.Response(context, base.USE_OR_PASSSWORD_ERROR, nil)
		return
	}
	db := context.MustGet("db").(*gorm.DB)
	dbAccount := service.FindByName(db, account.UserName)
	if 0 >= dbAccount.Id {
		base.Response(context, base.USE_NOT_EXIST, nil)
		return
	}
	mdString := strings.Trim(account.PassWord, "") + dbAccount.Salt
	md5Password := util.EncodeMd5(mdString)
	if md5Password != dbAccount.PassWord {
		base.Response(context, base.USE_OR_PASSSWORD_ERROR, nil)
		return
	}
	if token, ok := utils.GenerateToken(account.UserName, dbAccount.Id); nil != ok {
		base.Response(context, base.SERVER_ERROR, nil)
		return
	} else {
		dbAccount.Token = token
		entity.Clear(&dbAccount)
		base.Response(context, base.SUCCESS, dbAccount)
		return
	}
}

func Register(context *gin.Context) {
	var account entity.User
	if err := context.ShouldBindJSON(&account); err != nil {
		base.Response(context, base.NOT_LOGIN, nil)
		return
	}

	if "" == strings.Trim(account.UserName, "") || "" == strings.Trim(account.PassWord, "") {
		base.ResponseWithMessage(context, base.NOT_LOGIN, "用户名和密码不能为空", nil)
		return
	}

	db := context.MustGet("db").(*gorm.DB)
	dbAccount := service.FindByName(db, account.UserName)
	if dbAccount.Id > 0 {
		base.ResponseWithMessage(context, base.NOT_LOGIN, "用户名已经存在", nil)
		return
	}
	u1 := uuid.Must(uuid.NewV4())
	account.Salt = u1.String()
	account.PassWord = util.EncodeMd5(strings.Trim(account.PassWord, "") + dbAccount.Salt)
	if service.Add(db, &account) != nil {
		base.Response(context, base.SERVER_ERROR, nil)
		return
	} else {
		base.Response(context, base.SUCCESS, account)
		return
	}
}
