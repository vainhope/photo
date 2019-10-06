package entity

import (
	"database/sql/driver"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type TimeStamp struct {
	time.Time
}

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	//格式化秒
	seconds := t.Unix()
	if seconds < 0 {
		//为负数
		return []byte("0"),nil
	}
	return []byte(strconv.FormatInt(seconds, 10)), nil
}
func (t TimeStamp) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
func (t *TimeStamp) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = TimeStamp{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

type User struct {
	Id            int64     `gorm:"column:id" json:"id"`
	UserName      string    `gorm:"column:userName" json:"userName"`
	PassWord      string    `gorm:"column:password" json:"password"`
	Salt          string    `gorm:"column:salt" json:"salt,omitempty"`
	Nickname      string    `gorm:"column:nickname" json:"nickname"`
	Avatar        string    `gorm:"column:avatar" json:"avatar"`
	Sex           int       `gorm:"column:sex" json:"sex"`
	Type          int       `gorm:"column:type" json:"type"`
	Deleted       int       `gorm:"column:deleted" json:"deleted"`
	State         int       `gorm:"column:state" json:"state"`
	Phone         string    `gorm:"column:phone" json:"phone"`
	Email         string    `gorm:"column:email" json:"email"`
	Signature     string    `gorm:"column:signature" json:"signature"`
	Birthday      TimeStamp `gorm:"column:birthday"  json:"birthday"`
	CreateTime    TimeStamp `gorm:"column:createTime" json:"createTime"`
	ModifyTime    TimeStamp `gorm:"column:modifyTime"  json:"modifyTime"`
	LastLoginTime TimeStamp `gorm:"column:lastLoginTime"  json:"lastLoginTime"`
	Token         string    `gorm:"-" json:"token"`
}

//指定对应表别名
func (user User) TableName() string {
	return "t_user"
}

//执行后的钩子
func (u *User) AfterFind(db *gorm.DB) (err error) {
	db.Model(u).UpdateColumn("lastLoginTime", time.Now())
	return
}

func (u *User) BeforeSave(db *gorm.DB) {
	u.CreateTime = TimeStamp{time.Now()}
	if 0 == u.Type {
		u.Type = 1
	}
	return
}

func Clear(u *User) {
	u.Salt = ""
	u.PassWord = ""
}
