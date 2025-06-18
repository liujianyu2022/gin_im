package model

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FormId   uint   // 发送者
	TargetId uint   // 接收者
	Type     int // 消息传播类型，比如：群聊、私聊、广播等
	Media    int    // 消息内容类型，比如：文字、图片、音频等
	Content  string // 消息内容

}

func (table *Message) TableName() string {
	return "message"
}

// type Node struct {
// 	Conn      *websocket.Conn
// 	DataQueue chan []byte
// 	GroupSets set.Interface
// }

// // 映射关系
// var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// // 读写锁
// var rwLocker sync.RWMutex

// func Chat(writer http.ResponseWriter, request *http.Request) {
// 	query := request.URL.Query()

// 	userId := query.Get("userId")
// 	userIdInt64, _ := strconv.ParseInt(userId, 10, 64)
// 	// targetId := query.Get("targetId")
// 	// context := query.Get("context")
// 	// messageType := query.Get("type")

// 	conn, err := (&websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}).Upgrade(writer, request, nil)

// 	if err != nil {
// 		fmt.Println("conn error", err)
// 		return
// 	}

// 	// 获取conn
// 	node := &Node{
// 		Conn:      conn,
// 		DataQueue: make(chan []byte, 50),
// 		GroupSets: set.New(set.ThreadSafe),
// 	}

// 	// 获取用户关系

// 	// userId 和 node 绑定
// 	rwLocker.Lock()
// 	clientMap[int64(userIdInt64)] = node
// 	rwLocker.Unlock()

// 	// 发送逻辑
// 	go sendProc(node)

// 	// 接收逻辑
// 	go receiveProc(node)

// 	sendMessage(userIdInt64, []byte("welcome"))
// }

// func sendProc(node *Node) {
// 	for {
// 		select {
// 		case data := <-node.DataQueue:
// 			err := node.Conn.WriteMessage(websocket.TextMessage, data)
// 			if err != nil {
// 				fmt.Println("conn error", err)
// 				return
// 			}
// 		}
// 	}
// }

// func receiveProc(node *Node) {
// 	for {

// 		_, data, err := node.Conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("conn error", err)
// 			return
// 		}
// 		broadMessage(data)
// 		fmt.Println("[WS] <<<< ", data)
// 	}
// }

// var spendChan chan []byte = make(chan []byte, 1024)

// func broadMessage(data []byte) {
// 	spendChan <- data
// }

// func init() {
// 	go udpSendProc()

// 	go udpReceiveProc()
// }

// // 完成 udp 数据发送协程
// func udpSendProc() {
// 	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
// 		IP:   net.IPv4(192, 168, 79, 1),
// 		Port: 8080,
// 	})
// 	defer con.Close()

// 	if err != nil {
// 		fmt.Println("udpSendProc", err)
// 		return
// 	}

// 	for {
// 		select {
// 		case data := <-spendChan:
// 			_, err := con.Write(data)
// 			if err != nil {
// 				fmt.Println("udpSendProc", err)
// 				return
// 			}
// 		}
// 	}
// }

// // 完成 udp 数据接收协程
// func udpReceiveProc() {
// 	con, err := net.ListenUDP("udp", &net.UDPAddr{
// 		IP:   net.IPv4zero,
// 		Port: 8080,
// 	})

// 	defer con.Close()

// 	if err != nil {
// 		fmt.Println("udpReceiveProc", err)
// 		return
// 	}

// 	for {
// 		var buf [512]byte
// 		n, err := con.Read(buf[0:])

// 		if err != nil {
// 			fmt.Println("udpReceiveProc", err)
// 			return
// 		}

// 		dispatch(buf[0:n])
// 	}
// }

// // 后端调度
// func dispatch(data []byte) {
// 	msg := Message{}

// 	err := json.Unmarshal(data, &msg)
// 	if err != nil {
// 		fmt.Println("dispatch", err)
// 		return
// 	}

// 	switch msg.Type {
// 	case 1:
// 		sendMessage( int64(msg.TargetId), data)			// 私信
// 	// case 2:
// 	// 	sendGroupMessage()		// 群发
// 	// case 3:	
// 	// 	sendBroadMessage()		// 广播
// 	}
// }

// func sendMessage(targetId int64, msg []byte){
// // userId 和 node 绑定
// 	rwLocker.RLock()
// 	node, ok := clientMap[targetId]
// 	rwLocker.RUnlock()

// 	if ok {
// 		node.DataQueue <- msg
// 	}
// }
