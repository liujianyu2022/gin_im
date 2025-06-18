package service

import (
	"encoding/json"
	"fmt"
	"gin_im/dto"
	"gin_im/repository"

	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	Repository      *repository.RedisRepository
	ClientMap map[int64]*dto.WSNode
	RwLocker  sync.RWMutex
}

func NewWebsocketService(repository *repository.RedisRepository) *WebsocketService {
	return &WebsocketService{
		Repository: repository,
		ClientMap: make(map[int64]*dto.WSNode),
	}
}

// HandleConnection 处理WebSocket连接的核心逻辑
func (s *WebsocketService) Connection(node *dto.WSNode, userID int64) {
	// 1. 保存连接
	s.RwLocker.Lock()
	s.ClientMap[userID] = node
	s.RwLocker.Unlock()

	// 2. 启动读写协程
	go s.sendProc(node)
	go s.receiveProc(node, userID)
}

func (s *WebsocketService) sendProc(node *dto.WSNode) {
	defer func() {
		if err := node.Conn.Close(); err != nil {
			fmt.Println("WebSocket close error:", err)
		}
	}()

	for data := range node.DataQueue {
		if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("WebSocket write error:", err)
			return
		}
	}
}

func (s *WebsocketService) receiveProc(node *dto.WSNode, userID int64) {
	defer s.cleanupConnection(userID)

	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("[ws] <<< ", data)

		var msg dto.WSMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case 1: // 私聊
			s.sendToUser(msg.TargetId, data)
			// 可扩展其他消息类型
		}
	}
}

func (s *WebsocketService) sendToUser(targetID int64, msg []byte) {
	s.RwLocker.RLock()
	node, ok := s.ClientMap[targetID]
	s.RwLocker.RUnlock()

	if ok {
		node.DataQueue <- msg
	}
}

func (s *WebsocketService) cleanupConnection(userID int64) {
	s.RwLocker.Lock()
	defer s.RwLocker.Unlock()

	if node, exists := s.ClientMap[userID]; exists {
		close(node.DataQueue)
		delete(s.ClientMap, userID)
	}
}

// func (s *WebsocketService) HandleConnection(c *gin.Context) {
// 	// 获取用户ID
// 	userIDStr := c.Query("userId")
// 	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

// 	// 升级为WebSocket连接
// 	var upgrader = websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

// 	if err != nil {
// 		return
// 	}

// 	// 创建节点
// 	node := &dto.WSNode{
// 		Conn:      conn,
// 		DataQueue: make(chan []byte, 50),
// 	}

// 	// 保存连接
// 	s.RwLocker.Lock()
// 	s.ClientMap[userID] = node
// 	s.RwLocker.Unlock()

// 	// 启动读写协程
// 	go s.sendProc(node)
// 	go s.receiveProc(node, userID)
// }

// func (s *WebsocketService) sendProc(node *dto.WSNode) {
// 	defer func() {
// 		if err := node.Conn.Close(); err != nil {
// 			fmt.Println("WebSocket close error:", err)
// 		}
// 	}()

// 	for data := range node.DataQueue {
// 		if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
// 			fmt.Println("WebSocket write error:", err)
// 			return
// 		}
// 	}
// }

// func (s *WebsocketService) receiveProc(node *dto.WSNode, userID int64) {
// 	for {
// 		_, data, err := node.Conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		fmt.Println("[WS] <<<< ", data)

// 		var msg dto.WSMessage
// 		if err := json.Unmarshal(data, &msg); err != nil {
// 			continue
// 		}

// 		switch msg.Type {
// 		case 1: // 私聊
// 			s.sendToUser(msg.TargetId, data)
// 			// 可以扩展其他消息类型
// 		}
// 	}
// }

// func (s *WebsocketService) sendToUser(targetID int64, msg []byte) {
// 	s.RwLocker.RLock()
// 	node, ok := s.ClientMap[targetID]
// 	s.RwLocker.RUnlock()

// 	if ok {
// 		node.DataQueue <- msg
// 	}
// }
