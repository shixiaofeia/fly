package model

type (
	SocketHandleReq struct {
		OpType uint   `json:"opType"`
		Data   string `json:"data"`
	}
)
