package ws

import (
	"fmt"
)

// addUser 添加用户
func (c *SocketConn) addUser() {
	if c.UserId == "" {
		return
	}
	sockets.userMu.Lock()
	defer sockets.userMu.Unlock()
	sockets.Users[c.UserId] = c.ConnId
}

// delUser 删除用户
func (c *SocketConn) delUser() {
	if c.UserId == "" {
		return
	}
	sockets.userMu.Lock()
	defer sockets.userMu.Unlock()
	delete(sockets.Users, c.UserId)
}

// GetConnByUserId 获取指定用户连接
func (c *SocketConn) GetConnByUserId(userId UserId) (*SocketConn, error) {
	sockets.userMu.RLock()
	defer sockets.userMu.RUnlock()
	connId, ok := sockets.Users[userId]
	if !ok {
		return nil, fmt.Errorf("userId not exist")
	}
	return c.GetConnById(connId)
}
