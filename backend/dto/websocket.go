package dto

import "github.com/gorilla/websocket"

// 在服务端维护每个 WebSocket 客户端的连接状态和信息。
// WebsocketNode 是每个连接一个实例
type WebsocketNode struct {
	Conn      *websocket.Conn // WebSocket 连接对象。通过它可以进行消息的读取和写入操作
	DataQueue chan []byte     // 数据消息通道。一个缓冲通道，用于存储待发送给该客户端的数据
}

type WebsocketMessage struct {
	FromId    uint  `json:"from_id"`
	TargetId  uint  `json:"target_id"`			   // 对于私聊，则是对方用户ID；对于群聊，这里是群ID
	Type      int    `json:"type"`
	Content   string `json:"content"`
	GroupName string `json:"group_name,omitempty"` // 群名称(可选)
}


// chan 是一种通道类型，用于在不同的 Goroutine（协程）之间安全地传递数据，是实现并发编程中同步和通信的核心机制
// chan 是 Go 的内置类型，属于引用类型（底层通过指针实现）,遵循先进先出（FIFO）的规则，确保数据按发送顺序接收
// 阻塞机制：发送数据到通道时，如果通道已满（有缓冲且满），发送方会阻塞，直到数据被接收。从通道接收数据时，如果通道为空，接收方会阻塞，直到有数据可读
// 方向性：通道可以声明为单向（只读 <-chan 或只写 chan<-），通常用于函数参数限制权限


