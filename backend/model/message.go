package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FormId   uint   // 发送者
	TargetId uint   // 接收者
	Type     string // 消息传播类型，比如：群聊、私聊、广播等
	Media    int    // 消息内容类型，比如：文字、图片、音频等
	Content  string // 消息内容

}

func (table *Message) TableName() string {
	return "message"
}
