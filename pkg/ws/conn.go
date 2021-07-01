package ws

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

// GetConnById 获取指定连接
func (c *SocketConn) GetConnById(connId ConnId) (*SocketConn, error) {
	sockets.ClientLock.RLock()
	defer sockets.ClientLock.RUnlock()
	if v, ok := sockets.Clients[connId]; ok {
		return v, nil
	}
	return nil, errors.New("connId not exist")
}

// Close 关闭连接
func (c *SocketConn) Close() {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if err := c.Conn.Close(); err != nil {
		return
	}
	c.closeCh <- 1
	c.closeCh <- 1
	c.delUser()
	// 退出群组
	for groupId := range c.Groups {
		NewGroup(groupId).Exit(c.ConnId)
	}
	sockets.ClientLock.Lock()
	defer sockets.ClientLock.Unlock()
	delete(sockets.Clients, c.ConnId)
	close(c.sendCh)
	close(c.closeCh)
	return
}

// JoinGroup 加入组
func (c *SocketConn) JoinGroup(groupId GroupId) error {
	if !new(SocketGroup).Exist(groupId) {
		return errors.New("groupId not exist")
	}
	NewGroup(groupId).Join(c.ConnId)
	return nil
}

// ExitGroup 退出组
func (c *SocketConn) ExitGroup(groupId GroupId) {
	NewGroup(groupId).Exit(c.ConnId)
	c.Lock.Lock()
	defer c.Lock.Unlock()
	delete(c.Groups, groupId)
	return
}

// SendMsg 发送消息
func (c *SocketConn) SendMsg(msg interface{}) {
	b, _ := json.Marshal(msg)
	c.sendCh <- b
}

// production 生产
func (c *SocketConn) production() {
	for {
		select {
		case data := <-c.sendCh:
			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			}
		case <-c.closeCh:
			return
		}
	}
}

// consumer 消费
func (c *SocketConn) consumer(handle func(*SocketConn, []byte)) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			// 预防客户端未发送消息主动断开连接
			select {
			case <-c.closeCh:
				return
			case <-time.After(1 * time.Second):
				go c.Close()
				continue
			}
		}
		handle(c, data)
	}
}
