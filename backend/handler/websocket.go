package handler

import (
	"fmt"
	"gin_im/api"
	"gin_im/dto"
	"gin_im/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	Service *service.WebsocketService
}

func NewWebsocketHandler(websocketService *service.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{
		Service: websocketService,
	}
}

// 升级为WebSocket连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (handler *WebsocketHandler) Connect(ctx *gin.Context) {
	// 1. 协议层处理

	userIDStr := ctx.Query("userId")
	fmt.Println("userIDStr = ", userIDStr)

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Println("strconv.ParseInt(userIDStr, 10, 64):", err)
		api.HandleError(ctx, api.ErrUnauthorized, "")
		// return
	}

	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}

	// 创建WebsocketNode节点(包含WebSocket连接和数据通道)
	node := &dto.WebsocketNode{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
	}

	// 调用Service层的Connection方法处理连接
	handler.Service.Connection(node, userID)
}


// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func (handler *WebsocketHandler) Connect(ctx *gin.Context) {
// 	// 升级为 WebSocket 连接
// 	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	// 获取用户ID (从JWT中)
// 	userID, exists := ctx.Get("userID")
// 	if !exists {
// 		return
// 	}

// 	// 处理 WebSocket 连接
// 	handler.Service.HandleConnection(ctx, conn, userID.(string))
// }

// func  (handler *WebsocketHandler) Test(ctx *gin.Context) {
// 	handler.Service.SendUserMessage(ctx)
// }
