package ws

// NewGroup 获取组.
func NewGroup(groupId GroupId) *SocketGroup {
	sockets.groupMu.Lock()
	defer sockets.groupMu.Unlock()
	if g, ok := sockets.Groups[groupId]; ok {
		return g
	}
	sockets.Groups[groupId] = &SocketGroup{
		GroupId: groupId,
		ConnIds: make(map[ConnId]struct{}),
	}
	return sockets.Groups[groupId]
}

// Exist 组是否存在.
func (g *SocketGroup) Exist(groupId GroupId) bool {
	sockets.groupMu.RLock()
	defer sockets.groupMu.RUnlock()
	if _, ok := sockets.Groups[groupId]; ok {
		return true
	}
	return false
}

// Join 加入组.
func (g *SocketGroup) Join(connId ConnId) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.ConnIds[connId] = struct{}{}
}

// Exit 退出组.
func (g *SocketGroup) Exit(connId ConnId) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.ConnIds, connId)
}

// SendMsg 组内发送消息.
func (g *SocketGroup) SendMsg(msg interface{}) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	s := &SocketConn{}
	for connId := range g.ConnIds {
		if c, err := s.GetConnById(connId); err == nil {
			c.SendMsg(msg)
		}
	}
}
