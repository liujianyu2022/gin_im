package dto

import "github.com/gorilla/websocket"

type WSNode struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
}

type WSMessage struct {
	FormId   int64  `json:"form_id"`
	TargetId int64  `json:"target_id"`
	Type     int    `json:"type"`
	Content  string `json:"content"`
}
