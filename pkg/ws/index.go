package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type (
	ConnId  string // 连接唯一标识
	GroupId string // 组唯一标识
	UserId  string // 用户Id
	// Socket 连接管理
	Socket struct {
		ClientLock sync.RWMutex
		GroupLock  sync.RWMutex
		UserLock   sync.RWMutex
		Clients    map[ConnId]*SocketConn   // 连接Map
		Groups     map[GroupId]*SocketGroup // 组Map
		Users      map[UserId]ConnId        // 用户映射
	}
	// SocketGroup 组
	SocketGroup struct {
		Lock    sync.RWMutex
		GroupId GroupId
		ConnIds map[ConnId]struct{}
	}
	// SocketConn 单个连接
	SocketConn struct {
		Lock    sync.Mutex
		ConnId  ConnId
		Conn    *websocket.Conn
		UserId  UserId
		Groups  map[GroupId]struct{} // 加入的群组
		sendCh  chan []byte          // 发送消息队列
		closeCh chan uint8           // 关闭通道
	}
)

var (
	sockets = &Socket{
		Clients: make(map[ConnId]*SocketConn),
		Groups:  make(map[GroupId]*SocketGroup),
		Users:   make(map[UserId]ConnId),
	}
)

func init() {
	//go printSockets()
}

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
	sockets.ClientLock.Lock()
	sockets.Clients[client.ConnId] = client
	sockets.ClientLock.Unlock()
	client.addUser()
	go client.consumer(handle)
	go client.production()
	return
}

// printSockets 打印sockets信息(调试堆栈释放)
func printSockets() {
	var (
		clientNum    int
		groupNum     int
		userNum      int
		groupUserNum int
	)
	for {
		sockets.ClientLock.RLock()
		clientNum = len(sockets.Clients)
		sockets.ClientLock.RUnlock()
		sockets.GroupLock.RLock()
		groupNum = len(sockets.Groups)
		if val, ok := sockets.Groups["10010"]; ok {
			val.Lock.RLock()
			groupUserNum = len(val.ConnIds)
			val.Lock.RUnlock()
		}
		sockets.GroupLock.RUnlock()
		sockets.UserLock.RLock()
		userNum = len(sockets.Users)
		sockets.UserLock.RUnlock()
		log.Printf("clientNum: %d, groupNum: %d, userNum: %d, groupUserNum: %d", clientNum, groupNum, userNum, groupUserNum)
		time.Sleep(10 * time.Second)
	}

}
