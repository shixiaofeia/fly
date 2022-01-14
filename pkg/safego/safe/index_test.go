package safe

import (
	"fmt"
	"testing"
	"time"
)

func TestSafeGo(t *testing.T) {
	Go(func() {
		var (
			total = 10
			size  = 0
		)
		fmt.Println(total / size)
	})
	time.Sleep(3 * time.Second)
	t.Log("end")
}
