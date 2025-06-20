package dto

import "github.com/gorilla/websocket"

// 在服务端维护每个 WebSocket 客户端的连接状态和信息。
type WebsocketNode struct {
	Conn      *websocket.Conn // WebSocket 连接对象。通过它可以进行消息的读取和写入操作
	DataQueue chan []byte     // 数据消息通道。一个缓冲通道，用于存储待发送给该客户端的数据
}

type WebsocketMessage struct {
	FromId    int64  `json:"from_id"`
	TargetId  int64  `json:"target_id"`			   // 对于私聊，则是对方用户ID；对于群聊，这里是群ID
	Type      int    `json:"type"`
	Content   string `json:"content"`
	GroupName string `json:"group_name,omitempty"` // 群名称(可选)
}
