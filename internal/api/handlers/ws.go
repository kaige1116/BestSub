package handlers

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/api/middleware"
	"github.com/bestruirui/bestsub/internal/api/router"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// wsHandler WebSocket处理器
type wsHandler struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]*Client
	mu       sync.RWMutex
}

// Client WebSocket客户端信息
type Client struct {
	conn       *websocket.Conn
	nameFilter string
	send       chan log.LogEntry
	mu         sync.RWMutex
}

// init 函数用于自动注册路由
func init() {
	h := newWSHandler()

	router.NewGroupRouter("/api/v1/ws").
		Use(middleware.WSAuth()).
		AddRoute(
			router.NewRoute("/logs", router.GET).
				Handle(h.handleLogWebSocket).
				WithDescription("WebSocket logs streaming"),
		)
}

// newWSHandler 创建WebSocket处理器
func newWSHandler() *wsHandler {
	h := &wsHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				if utils.IsDebug() {
					return true
				}

				origin := r.Header.Get("Origin")

				if origin == "" {
					log.Debugf("WebSocket客户端连接: 没有Origin头")
					return true
				}
				// TODO: 添加允许的域名列表

				log.Debugf("WebSocket客户端连接: Origin=%s", origin)

				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients: make(map[*websocket.Conn]*Client),
	}
	go h.broadcastLogs()
	return h
}

func (h *wsHandler) handleLogWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("WebSocket升级失败: %v", err)
		return
	}

	nameFilter := c.Query("name")

	client := &Client{
		conn:       conn,
		nameFilter: nameFilter,
		send:       make(chan log.LogEntry, 256),
	}

	h.mu.Lock()
	h.clients[conn] = client
	h.mu.Unlock()

	username, _ := c.Get("username")
	clientIP := c.ClientIP()
	log.Infof("WebSocket客户端连接: 用户=%s, IP=%s, 过滤器=%s", username, clientIP, nameFilter)

	go h.handleClient(client)
	go h.readPump(client)
}

func (h *wsHandler) broadcastLogs() {
	logChannel := log.GetWSChannel()

	for logEntry := range logChannel {
		h.mu.RLock()
		for _, client := range h.clients {
			if h.shouldSendLog(client, logEntry) {
				select {
				case client.send <- logEntry:
				default:
					close(client.send)
					delete(h.clients, client.conn)
				}
			}
		}
		h.mu.RUnlock()
	}
}

func (h *wsHandler) shouldSendLog(client *Client, logEntry log.LogEntry) bool {
	client.mu.RLock()
	defer client.mu.RUnlock()

	if client.nameFilter == "" {
		return true
	}

	return strings.Contains(logEntry.Name, client.nameFilter)
}

func (h *wsHandler) handleClient(client *Client) {
	defer func() {
		h.removeClient(client)
		client.conn.Close()
	}()

	client.conn.SetWriteDeadline(time.Now().Add(60 * time.Second))

	for logEntry := range client.send {
		client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		if err := client.conn.WriteJSON(logEntry); err != nil {
			log.Errorf("WebSocket发送消息失败: %v", err)
			return
		}
	}

	client.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (h *wsHandler) readPump(client *Client) {
	defer func() {
		h.removeClient(client)
		client.conn.Close()
	}()

	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("WebSocket连接异常关闭: %v", err)
			}
			break
		}
	}
}

func (h *wsHandler) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client.conn]; exists {
		delete(h.clients, client.conn)
		close(client.send)
		log.Debugf("WebSocket客户端断开连接")
	}
}
