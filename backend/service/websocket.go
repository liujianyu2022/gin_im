package service

import (
	"encoding/json"
	"fmt"
	"gin_im/dto"
	"gin_im/repository"

	"sync"

	"github.com/gorilla/websocket"
)

/*
假设有4个用户在线，他们的用户ID分别是1、2、3、4：
ClientMap = {
    1: &WebsocketNode{Conn: conn1, DataQueue: make(chan []byte, 50)},
    2: &WebsocketNode{Conn: conn2, DataQueue: make(chan []byte, 50)},
    3: &WebsocketNode{Conn: conn3, DataQueue: make(chan []byte, 50)},
    4: &WebsocketNode{Conn: conn4, DataQueue: make(chan []byte, 50)}
}

操作1：用户1发送私聊消息给用户3
a. 用户1的消息通过自己的连接到达服务端			用户1的 receiveProcedure 通过 conn1 收到消息
b. 服务端解析消息发现 TargetId=3			  服务端解析出 TargetId=3
c. 从 ClientMap 查找键为3的 WebsocketNode		      使用读锁查找 ClientMap[3] 获取用户3的 WebsocketNode
d. 将消息放入用户3的 DataQueue				   将消息放入用户3的 DataQueue，用户3的 sendProcedure 从 DataQueue 取出消息并通过 conn3 发送

操作2：用户2下线
连接断开触发 cleanupConnection
ClientMap = {
    1: &WebsocketNode{Conn: conn1, DataQueue: make(chan []byte, 50)},
    3: &WebsocketNode{Conn: conn3, DataQueue: make(chan []byte, 50)},
    4: &WebsocketNode{Conn: conn4, DataQueue: make(chan []byte, 50)}
}

操作3：新用户5连接
ClientMap = {
    1: &WebsocketNode{...},
    3: &WebsocketNode{...},
    4: &WebsocketNode{...},
    5: &WebsocketNode{Conn: conn5, DataQueue: make(chan []byte, 50)}
}
*/

type WebsocketService struct {
	Repository *repository.RedisRepository
	ClientMap  map[uint]*dto.WebsocketNode // 存储用户ID到 WebSocket 连接节点(WebsocketNode) 的映射
	RwLocker   sync.RWMutex                 // 读写锁，保护ClientMap的并发访问
}

func NewWebsocketService(repository *repository.RedisRepository) *WebsocketService {
	return &WebsocketService{
		Repository: repository,
		ClientMap:  make(map[uint]*dto.WebsocketNode),
	}
}

// HandleConnection 处理WebSocket连接的核心逻辑
// 将新连接存入ClientMap
// 启动两个goroutine分别处理发送和接收
func (s *WebsocketService) Connection(node *dto.WebsocketNode, userID uint) {

	// 1. 保存连接
	s.RwLocker.Lock() // 写操作(添加/删除连接)使用写锁 Lock()
	s.ClientMap[userID] = node
	s.RwLocker.Unlock()

	// 2. 启动读写协程
	// 在WebSocket通信中，"发送"和"接收"是相对于服务端而言的
	go s.receiveProcedure(node, userID) // 负责接收客户端发来的消息（客户端→服务端）
	go s.sendProcedure(node)            // 负责发送消息给客户端（服务端→客户端）
}

// 持续读取WebSocket消息
// 解析消息类型
// 根据消息类型调用相应处理方法(如sendToUser)
func (service *WebsocketService) receiveProcedure(node *dto.WebsocketNode, userID uint) {
	defer service.cleanupConnection(userID)					// 唯一关闭入口

	for {
		_, data, err := node.Conn.ReadMessage()

		if err != nil {
			break
		}

		fmt.Println("[ws] <<< ", data)

		var msg dto.WebsocketMessage

		if err := json.Unmarshal(data, &msg); err != nil {
			// 返回错误消息给客户端
            node.DataQueue <- []byte(`{"type":"error", "message":"invalid message format"}`)
			continue
		}

		switch msg.Type {
		case 1: // 私聊
			service.sendToUser(msg.TargetId, data)
		case 2: // 群聊
			service.sendToGroup(msg.TargetId, msg.FromId, data)
		}
	}
}

// 从 WebsocketNode节点 的DataQueue通道读取数据并发送到WebSocket连接
// 仅负责发送数据，不关心连接状态。
func (s *WebsocketService) sendProcedure(node *dto.WebsocketNode) {
	for data := range node.DataQueue {
		if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("WebSocket write error:", err)
			return
		}
	}
}

// 查找目标用户连接
// 将消息放入目标用户的数据队列
func (service *WebsocketService) sendToUser(targetID uint, msg []byte) {
	service.RwLocker.RLock() 									// 读操作(查找连接)使用读锁(RLock())
	node, ok := service.ClientMap[targetID]
	service.RwLocker.RUnlock()

	if ok {
		node.DataQueue <- msg
	}
}

// 发送消息到群组
func (service *WebsocketService) sendToGroup(groupId uint, fromId uint, msg []byte) {
	// 1. 查询群组成员列表
	members, err := service.Repository.GetGroupMembers(groupId)

	if err != nil {
		fmt.Printf("获取群组成员失败: %v\n", err)
		return
	}

	// 2. 遍历成员发送消息
	service.RwLocker.RLock()
	defer service.RwLocker.RUnlock()

	for _, memberId := range members {
		if node, ok := service.ClientMap[memberId]; ok {
			// 排除发送者自己
			if memberId == fromId {
				continue
			}

			select {
			case node.DataQueue <- msg:
				// 成功放入队列
			default:
				fmt.Printf("用户%d的消息队列已满，丢弃群消息\n", memberId)
			}
		}
	}
}

// 清理断开连接的资源(关闭通道、从ClientMap移除)
// 只需要调用一次就好了
func (service *WebsocketService) cleanupConnection(userID uint) {
	service.RwLocker.Lock()
	defer service.RwLocker.Unlock()

	node, exists := service.ClientMap[userID]

	if !exists {
        return
    }

    // 避免重复关闭
    if node.DataQueue != nil {
        close(node.DataQueue)
    }

    if node.Conn != nil {
        node.Conn.Close()
    }

    delete(service.ClientMap, userID)
}
