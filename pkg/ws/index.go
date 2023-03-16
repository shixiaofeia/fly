package ws

import (
	"fly/pkg/safego/safe"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

var (
	sockets = &Socket{
		Clients: make(map[ConnId]*SocketConn),
		Groups:  make(map[GroupId]*SocketGroup),
		Users:   make(map[UserId]ConnId),
	}
)

// NewClient 新的连接.
func NewClient(connId ConnId, userId UserId, conn *websocket.Conn, handle func(*SocketConn, []byte)) {
	client := &SocketConn{
		ConnId:  connId,
		Conn:    conn,
		UserId:  userId,
		Groups:  make(map[GroupId]struct{}),
		sendCh:  make(chan []byte, 20),
		closeCh: make(chan struct{}),
	}
	sockets.clientMu.Lock()
	sockets.Clients[client.ConnId] = client
	sockets.clientMu.Unlock()
	client.addUser()
	safe.Go(func() {
		client.consumer(handle)
	})
	safe.Go(client.production)
	return
}

func PrintSocketLength() {
	for {
		time.Sleep(10 * time.Second)
		sockets.clientMu.RLock()
		fmt.Println(fmt.Sprintf("clinets: %d, groups: %d, users: %d", len(sockets.Clients), len(sockets.Groups), len(sockets.Users)))
		sockets.clientMu.RUnlock()
	}
}
