package redis

import (
	"testing"
	"time"
)

func TestInitStandAlone(t *testing.T) {
	_ = Init(Conf{Addr: "127.0.0.1:6379"})
	now := time.Now().Unix()

	type item struct {
		name   string
		key    string
		member interface{}
		want   int64
	}

	tests := []item{
		{name: "sAdd success", key: "test", member: now, want: 1},
		{name: "sAdd repeat", key: "test", member: now, want: 0},
	}

	cl := NewStandAloneClient()
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if val := cl.SAdd(v.key, v.member).Val(); val != v.want {
				t.Errorf("SAdd got = %v, want = %v", val, v.want)
			}
		})
	}
}
