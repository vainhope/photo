package photo

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"goweb/base"
	"goweb/entity"
	"goweb/util"
	"strconv"
	"time"
)

func CoverCaptures(context *gin.Context) {
	coverId, ok := strconv.Atoi(context.Query("coverId"))
	if ok != nil || coverId < 0 {
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "请选择的正确的封面", nil)
		return
	}
	db := context.MustGet("db").(*gorm.DB)
	var ret []entity.Photo
	//具体图片信息
	db.Raw("select * from t_photo where deleted = 0 and belongCover = ? order by sort , createTime  ", coverId).Scan(&ret)
	//封面详情
	var cover entity.Cover
	db.Raw("select * from t_cover where id = ? and deleted = 0", coverId).Scan(&cover)
	data := make(map[string]interface{})

	data["cover"] = cover
	data["photos"] = ret
	base.Response(context, base.SUCCESS, data)
	return
}

func AddCaptures(context *gin.Context) {
	param := make(map[string]interface{})
	ok := context.Bind(&param)

	if ok != nil {
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "参数错误", nil)
		return
	}

	coverId := util.Wrap(param["coverId"].(float64), 0)
	uploadId := util.Wrap(param["uploadId"].(float64), 0)

	if uploadId <= 0 {
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "上传信息不能为空", nil)
		return
	}
	if coverId <= 0 {
		base.ResponseWithMessage(context, base.PARAMETER_ERROR, "封面不能为空", nil)
		return
	}

	db := context.MustGet("db").(*gorm.DB)
	photo := entity.Photo{}
	log := entity.UploadLog{}
	log.Id = uploadId
	db.First(&log)
	photo.CreateTime = entity.TimeStamp{Time: time.Now()}
	photo.ModifyTime = entity.TimeStamp{Time: time.Now()}
	photo.ImageUrl = log.Url
	photo.OriginName = log.Name
	photo.BelongCover = coverId
	photo.CreateUserId = context.MustGet("userId").(int64)
	photo.UploadId = uploadId
	db.Create(&photo)
	base.Response(context, base.SUCCESS, photo)
	return
}
