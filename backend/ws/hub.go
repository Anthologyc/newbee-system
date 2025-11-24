package ws

import "sync"

type Hub struct {
	clients    map[*Client]bool
	// 🚀 新增：房间管理
	Rooms      map[string]*Room 
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				// 🚀 如果在房间里，先退出房间
				if client.RoomID != "" {
					if room, ok := h.Rooms[client.RoomID]; ok {
						room.Leave(client)
						// 如果房间没人了，销毁房间
						if len(room.Clients) == 0 {
							delete(h.Rooms, client.RoomID)
						}
					}
				}
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
		
		// ... broadcast case ...
		}
	}
}

// 🚀 创建房间
func (h *Hub) CreateRoom(id string, hostID uint) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, ok := h.Rooms[id]; !ok {
		room := NewRoom(id, h)
		h.Rooms[id] = room
		go room.Run()
		return room
	}
	return h.Rooms[id]
}

// 🚀 获取房间
func (h *Hub) GetRoom(id string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.Rooms[id]
}

// 🚀 获取房间列表 (简单版)
func (h *Hub) GetRoomList() []map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var list []map[string]interface{}
	for id, r := range h.Rooms {
		list = append(list, map[string]interface{}{
			"id": id,
			"count": len(r.Clients),
			"status": r.Status,
		})
	}
	return list
}

var GlobalHub *Hub // 全局变量

func Init() {
    GlobalHub = NewHub()
    go GlobalHub.Run()
}