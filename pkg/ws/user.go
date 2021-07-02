package ws

import "errors"

// addUser 添加用户
func (c *SocketConn) addUser() {
	if c.UserId == "" {
		return
	}
	sockets.UserLock.Lock()
	defer sockets.UserLock.Unlock()
	sockets.Users[c.UserId] = c.ConnId
}

// delUser 删除用户
func (c *SocketConn) delUser() {
	if c.UserId == "" {
		return
	}
	sockets.UserLock.Lock()
	defer sockets.UserLock.Unlock()
	delete(sockets.Users, c.UserId)
}

// GetConnByUserId 获取指定用户连接
func (c *SocketConn) GetConnByUserId(userId UserId) (*SocketConn, error) {
	sockets.UserLock.RLock()
	defer sockets.UserLock.RUnlock()
	connId, ok := sockets.Users[userId]
	if !ok {
		return nil, errors.New("userId not exist")
	}
	return c.GetConnById(connId)
}
