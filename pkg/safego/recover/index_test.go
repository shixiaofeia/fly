package recover

import (
	"fmt"
	"testing"
	"time"
)

func TestSafeGo(t *testing.T) {
	SafeGo(func() {
		var (
			total = 10
			size  = 0
		)
		fmt.Println(total / size)
	})
	time.Sleep(3 * time.Second)
	t.Log("end")
}
