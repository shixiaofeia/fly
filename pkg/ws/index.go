package ws

import (
	recover2 "fly/pkg/safego/recover"
	"github.com/gorilla/websocket"
)

var (
	sockets = &Socket{
		Clients: make(map[ConnId]*SocketConn),
		Groups:  make(map[GroupId]*SocketGroup),
		Users:   make(map[UserId]ConnId),
	}
)

// NewClient 新的连接
func NewClient(connId ConnId, userId UserId, conn *websocket.Conn, handle func(*SocketConn, []byte)) {
	client := &SocketConn{
		ConnId:  connId,
		Conn:    conn,
		UserId:  userId,
		Groups:  make(map[GroupId]struct{}),
		sendCh:  make(chan []byte, 20),
		closeCh: make(chan uint8, 2),
	}
	sockets.clientMu.Lock()
	sockets.Clients[client.ConnId] = client
	sockets.clientMu.Unlock()
	client.addUser()
	recover2.SafeGo(func() {
		client.consumer(handle)
	})
	recover2.SafeGo(client.production)
	return
}
