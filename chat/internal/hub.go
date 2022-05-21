package internal

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool // 上线clients

	// Inbound messages from the clients.
	broadcast chan []byte // 客户端发送的消息 ->广播给其他的客户端

	// Register requests from the clients.
	register chan *Client // 注册channel，接收注册msg

	// Unregister requests from clients.
	unregister chan *Client // 下线channel
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: // 注册channel：存放到注册表中，数据流也就在这些client中发生
			h.clients[client] = true
		case client := <-h.unregister: // 下线channel：从注册表里面删除
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast: // 广播消息：发送给注册表中的client中，send接收到并显示到client上
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
