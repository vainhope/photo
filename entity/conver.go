package entity

type Page struct {
	CurrentPage int `gorm:"-" json:"currentPage"`
	PageSize    int `gorm:"-" json:"pageSize"`
}

type Cover struct {
	Id           int64     `gorm:"column:id;primary_key" json:"id"`
	UploadId     int64     `gorm:"-" json:"uploadId"`
	Cover        string    `gorm:"column:cover" json:"cover"`
	Desc         string    `gorm:"column:desc" json:"desc"`
	Title        string    `gorm:"column:title" json:"title"`
	Private      int       `gorm:"column:private" json:"private"`
	Deleted      int       `gorm:"column:deleted" json:"deleted"`
	CreateTime   TimeStamp `gorm:"column:createTime" json:"createTime"`
	CreateUserId int64     `gorm:"column:createUserId" json:"createUserId"`
	ModifyTime   TimeStamp `gorm:"column:modifyTime"  json:"modifyTime"`
}

//指定对应表别名
func (cover Cover) TableName() string {
	return "t_cover"
}
