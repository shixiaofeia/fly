package safe

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
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

func TestGoWithField(t *testing.T) {
	var (
		wg       = new(sync.WaitGroup)
		goMap    = NewMap()
		fieldMap = NewMap()
		num      = 1000
	)

	for i := 0; i < num; i++ {
		wg.Add(2)
		Go(func() {
			goMap.Add(i, nil)
			wg.Done()
		})
		GoWithField(func(val interface{}) {
			fieldMap.Add(val, nil)
			wg.Done()
		}, i)
	}
	wg.Wait()
	assert.Equal(t, fieldMap.Len(), num)
	assert.Less(t, goMap.Len(), num)
}
