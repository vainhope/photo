package entity

import "time"

type UploadLog struct {
	Id         int64     `gorm:"column:id;primary_key" json:"id"`
	Name       string    `gorm:"column:Name" json:"name"`
	Length     int64     `gorm:"column:length" json:"length"`
	UserId     int64     `gorm:"column:userId" json:"userId"`
	Url        string    `gorm:"column:url" json:"url"`
	DeleteUrl  string    `gorm:"column:deleteUrl" json:"deleteUrl"`
	CreateTime time.Time `gorm:"column:createTime" json:"createTime"`
	Path       string    `gorm:"column:path" json:"path"`
}

//指定对应表别名
func (convers UploadLog) TableName() string {
	return "t_upload_log"
}
