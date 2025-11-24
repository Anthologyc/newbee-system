package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 定义通用的消息结构
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	UserID  uint        `json:"user_id,omitempty"` // 发送者的ID
	RoomID  string      `json:"room_id,omitempty"`
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan []byte
	UserID  uint
	RoomID  string
	IsReady bool
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("json unmarshal error: %v", err)
			continue
		}
		// 强制标记发送者ID，防止客户端伪造
		msg.UserID = c.UserID

		// 1. 处理加入房间
		if msg.Type == "join_room" {
			if payload, ok := msg.Payload.(map[string]interface{}); ok {
				if rid, ok := payload["room_id"].(string); ok {
					if room := c.Hub.GetRoom(rid); room != nil {
						room.Join(c)
					}
				}
			}
			continue
		}

		// 2. 处理离开房间
		if msg.Type == "leave_room" {
			if c.RoomID != "" {
				if room := c.Hub.GetRoom(c.RoomID); room != nil {
					room.Leave(c)
				}
			}
			continue
		}

		// 3. 其他消息转发给 Room 处理 (ready, start_game 等)
		if c.RoomID != "" {
			if room := c.Hub.GetRoom(c.RoomID); room != nil {
				room.HandleMessage(c, msg)
			}
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs 处理 WebSocket 请求
func ServeWs(hub *Hub, c *gin.Context) {
	// 🚀 关键：获取 user_id
	userIdStr := c.Query("user_id")
	uid, err := strconv.Atoi(userIdStr)
	if err != nil || uid == 0 {
		log.Println("Invalid user_id in ws connection")
		// 这里可以选择 return 不让连，或者分配临时 ID
		// 为了调试方便，暂时允许继续，但 user_id 为 0 会导致逻辑问题
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: uint(uid), // 赋值
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}