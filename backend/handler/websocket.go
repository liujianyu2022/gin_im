package handler

import (
	"encoding/json"
	"fmt"
	"gin_im/api"
	"gin_im/config"
	"gin_im/dto"
	"gin_im/service"
	"gin_im/tools"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	Service *service.WebsocketService
	Config  *config.Config
}

func NewWebsocketHandler(websocketService *service.WebsocketService, config *config.Config) *WebsocketHandler {
	return &WebsocketHandler{
		Service: websocketService,
		Config: config,
	}
}

// 升级为WebSocket连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (handler *WebsocketHandler) Connect(ctx *gin.Context) {

	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		api.HandleError(ctx, api.ErrInternalServer, "WebSocket upgrade error")
		return
	}

	// 创建WebsocketNode节点(包含WebSocket连接和数据通道)
	// 每个连接对应一个WebsocketNode实例，然后这些实例统一被 WebsocketService 的 ClientMap 所管理
	// 创建节点但不立即注册到ClientMap
	node := &dto.WebsocketNode{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
	}

	// 启动认证和处理流程
	go handler.handleConnection(ctx, node)
}

// handleConnection 处理认证相关逻辑
func (handler *WebsocketHandler) handleConnection(ctx *gin.Context, node *dto.WebsocketNode) {

	// 1. 认证阶段
	userID, err := handler.authenticate(node)
	if err != nil {
		api.HandleError(ctx, api.ErrBadRequest, err)
		return
	}

	// 2. 将连接交给Service管理
	handler.Service.Connection(node, userID)
	
}

// authenticate 处理WebSocket连接的认证
func (handler *WebsocketHandler) authenticate(node *dto.WebsocketNode) (uint, error) {

	// 读取认证消息
	_, authMsg, err := node.Conn.ReadMessage()

	if err != nil {
		return 0, fmt.Errorf("failed to read auth message: %v", err)
	}

	// 解析认证消息
	var auth struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	}

	if err := json.Unmarshal(authMsg, &auth); err != nil {
		node.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","message":"invalid auth format"}`))
		return 0, fmt.Errorf("invalid auth format: %v", err)
	}

	if auth.Type != "auth" {
		node.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","message":"first message must be auth"}`))
		return 0, fmt.Errorf("first message must be auth")
	}

	// 验证token
	fmt.Println("auth = ", auth)
	claims, err := tools.ParseToken(auth.Token, handler.Config)
	if err != nil {
		node.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","message":"invalid token"}`))
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	// 发送认证成功响应
	node.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_success", "user_id":` + strconv.FormatUint(uint64(claims.UserID), 10)+`}`))

	return claims.UserID, nil
}
