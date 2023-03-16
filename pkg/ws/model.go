package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type (
	ConnId  string // 连接唯一标识
	GroupId string // 组唯一标识
	UserId  string // 用户Id
	// Socket 连接管理
	Socket struct {
		clientMu sync.RWMutex
		groupMu  sync.RWMutex
		userMu   sync.RWMutex
		Clients  map[ConnId]*SocketConn   // 连接Map
		Groups   map[GroupId]*SocketGroup // 组Map
		Users    map[UserId]ConnId        // 用户映射
	}
	// SocketGroup 组
	SocketGroup struct {
		mu      sync.RWMutex
		GroupId GroupId
		ConnIds map[ConnId]struct{}
	}
	// SocketConn 单个连接
	SocketConn struct {
		mu      sync.Mutex
		ConnId  ConnId
		Conn    *websocket.Conn
		UserId  UserId
		Groups  map[GroupId]struct{} // 加入的群组
		sendCh  chan []byte          // 发送消息队列
		closeCh chan struct{}        // 关闭通道
	}
)
