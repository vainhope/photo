package cover

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"goweb/base"
	"goweb/entity"
	"goweb/service"
	"time"
)

func Add(context *gin.Context) {
	var cover entity.Cover
	if err := context.ShouldBindJSON(&cover); err != nil {
		base.Response(context, base.PARAMETER_ERROR, nil)
		return
	}
	if "" == cover.Title {
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "封面名不能为空", nil)
		return
	}
	if "" == cover.Cover {
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "封面不能为空", nil)
		return
	}
	cover.CreateUserId = context.MustGet("userId").(int64)
	db := context.MustGet("db").(*gorm.DB)
	cover.CreateTime = entity.TimeStamp{Time: time.Now()}
	cover.ModifyTime = entity.TimeStamp{Time: time.Now()}
	if db.Create(&cover) != nil {
		//添加photo
		_ = service.AddPhotoByUploadLog(db, cover.UploadId, cover.Id, cover.CreateUserId)
		base.Response(context, base.SUCCESS, cover)
	} else {
		base.Response(context, base.SUCCESS, cover)
	}
	return
}

func PageQuery(context *gin.Context) {
	var page entity.Page
	if err := context.ShouldBindJSON(&page); err != nil {
		base.Response(context, base.PARAMETER_ERROR, nil)
		return
	}
	if page.CurrentPage <= 0 {
		page.CurrentPage = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 20
	}

	db := context.MustGet("db").(*gorm.DB)
	var createUserId int64
	//TODO 需要解析token
	if id, ok := context.Get("userId"); !ok {
		createUserId = 0
	} else {
		createUserId = id.(int64)
	}
	var ret []entity.Cover
	db.Raw("SELECT * from t_cover where deleted = 0 and"+
		" (private = 0 or (private = 1 and createUserId = ? ))  ORDER by private,id DESC "+
		"limit ? OFFSET ?", createUserId, page.PageSize, (page.CurrentPage-1)*page.PageSize).Scan(&ret)
	base.ResponseWithPage(context, base.SUCCESS, ret, page)
	return
}
