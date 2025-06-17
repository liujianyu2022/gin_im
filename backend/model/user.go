package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]\\d{9}$)"` // 修正标签
	Email         string `valid:"email"`                   // 修正标签
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartbeatTime uint64
	LogoutTime    uint64
	IsLogout      bool
	DeviceInfo    string
}

func (table *User) TableName() string {
	return "user"
}
