package controller

import (
	"encoding/json"
	"fly/example/api/example/model"
	"fly/pkg/ws"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// SocketHealth socket.
func SocketHealth(ctx *gin.Context) {
	upgrade := websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	// 生产环境要保证Id唯一
	connId := uuid.NewV4().String()
	ws.NewClient(ws.ConnId(connId), ws.UserId(connId), conn, socketHandle)
}

// socketHandle socket处理.
func socketHandle(conn *ws.SocketConn, data []byte) {
	req := &model.SocketHandleReq{}
	_ = json.Unmarshal(data, req)
	switch req.OpType {
	case 1:
		conn.SendMsg("hello: " + req.Data)
	case 2:
		_ = conn.JoinGroup(ws.GroupId(req.Data))
	case 3:
		conn.Close()
	}
}

// GroupSend 测试组发送.
func GroupSend() {
	g := ws.NewGroup("10010")
	for {
		g.SendMsg("group send msg")
		time.Sleep(10 * time.Second)
	}
}
