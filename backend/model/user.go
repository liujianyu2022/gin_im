package model

import "gorm.io/gorm"

type User struct{
	gorm.Model
	Name string
	Password string
	Phone string
	Email string
	Identity string
	ClientIp string
	ClientPort string
	LoginTime uint64
	HeartbeatTime uint64
	LogoutTime uint64
	IsLogout bool
	DeviceInfo string
}

func (table *User) TableName() string {
	return "user"
}