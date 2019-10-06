package entity

type Photo struct {
	Id           int64     `gorm:"column:id" json:"id"`
	OriginName   string    `gorm:"column:originName" json:"originName"`
	UuidName     string    `gorm:"column:uuidName" json:"uuidName"`
	BelongCover  int64     `gorm:"column:belongCover" json:"belongCover"`
	CreateUserId int64     `gorm:"column:createUserId" json:"createUserId"`
	UploadId     int64     `gorm:"column:uploadId" json:"uploadId"`
	ImageUrl     string    `gorm:"column:imageUrl" json:"imageUrl"`
	CreateTime   TimeStamp `gorm:"column:createTime" json:"createTime"`
	ModifyTime   TimeStamp `gorm:"column:modifyTime"  json:"modifyTime"`
}

func (photo Photo) TableName() string {
	return "t_photo"
}
