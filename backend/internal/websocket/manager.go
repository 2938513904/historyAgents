package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	conn   *websocket.Conn
	roomID string
	send   chan []byte
}

// Manager WebSocket连接管理器
type Manager struct {
	clients    map[string]map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// BroadcastMessage 广播消息结构
type BroadcastMessage struct {
	RoomID  string
	Message []byte
}

// NewManager 创建新的WebSocket管理器
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Start 启动WebSocket管理器
func (m *Manager) Start() {
	go func() {
		for {
			select {
			case client := <-m.register:
				m.mutex.Lock()
				if m.clients[client.roomID] == nil {
					m.clients[client.roomID] = make(map[*Client]bool)
				}
				m.clients[client.roomID][client] = true
				m.mutex.Unlock()

			case client := <-m.unregister:
				m.mutex.Lock()
				if roomClients, ok := m.clients[client.roomID]; ok {
					if _, ok := roomClients[client]; ok {
						close(client.send)
						delete(roomClients, client)
						if len(roomClients) == 0 {
							delete(m.clients, client.roomID)
						}
					}
				}
				m.mutex.Unlock()

			case message := <-m.broadcast:
				m.mutex.RLock()
				roomClients, exists := m.clients[message.RoomID]
				m.mutex.RUnlock()

				if exists {
					for client := range roomClients {
						select {
						case client.send <- message.Message:
						default:
							close(client.send)
							m.mutex.Lock()
							delete(m.clients[message.RoomID], client)
							if len(m.clients[message.RoomID]) == 0 {
								delete(m.clients, message.RoomID)
							}
							m.mutex.Unlock()
						}
					}
				}
			}
		}
	}()
}

// BroadcastToRoom 向指定房间广播消息
func (m *Manager) BroadcastToRoom(roomID string, message []byte) {
	m.broadcast <- BroadcastMessage{
		RoomID:  roomID,
		Message: message,
	}
}

// GetClientsCount 获取指定房间的客户端数量
func (m *Manager) GetClientsCount(roomID string) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if roomClients, ok := m.clients[roomID]; ok {
		return len(roomClients)
	}
	return 0
}
