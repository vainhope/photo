package service

import (
	"github.com/jinzhu/gorm"
	"goweb/base"
	"goweb/entity"
	"time"
)

func AddPhotoByUploadLog(db *gorm.DB, logId int64, coverId int64, userId int64) error {
	if coverId < 0 {
		panic(base.Err(-1, "请选择目录"))
	}
	if logId < 0 {
		panic(base.Err(-1, "上传文件不存在"))
	}
	log := entity.UploadLog{Id: logId}
	db.First(&log)
	if log.Id < 0 {
		panic(base.Err(-1, "上传文件不存在"))
	}
	photo := entity.Photo{}
	photo.CreateTime = entity.TimeStamp{Time: time.Now()}
	photo.ModifyTime = entity.TimeStamp{Time: time.Now()}
	photo.ImageUrl = log.Url
	photo.OriginName = log.Name
	photo.BelongCover = coverId
	if 0 > userId {
		userId = 0
	}
	photo.CreateUserId = userId
	photo.UploadId = logId
	db.Create(&photo)
	return nil
}
