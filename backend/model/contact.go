package model

import "gorm.io/gorm"

// 人员关系表
type Contact struct {
	gorm.Model
	OwnerId uint
	TragetId uint
	Type int
	Description string
}

func (table *Contact) TableName() string {
	return "contact"
}
