package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string // 群聊的名称
	OwnerId     uint   // 群主
	Icon        string // 群图片
	Description string // 群聊描述
	Type        string // 预留字段
}

func (table *Group) TableName() string {
	return "group"
}