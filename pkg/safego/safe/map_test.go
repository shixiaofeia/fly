package safe

import (
	"sync"
	"testing"
)

func TestNewMap(t *testing.T) {
	var (
		wg = new(sync.WaitGroup)
		m  = NewMap()
	)

	for i := 0; i < 1000; i++ {
		wg.Add(5)

		GoWithField(func(val interface{}) {
			m.Add(val, nil)
			wg.Done()
		}, i)

		GoWithField(func(val interface{}) {
			m.Del(val)
			wg.Done()
		}, i)

		GoWithField(func(val interface{}) {
			m.Get(val)
			wg.Done()
		}, i)

		Go(func() {
			m.Len()
			wg.Done()
		})

		Go(func() {
			m.Range(func(k, v interface{}) bool {
				return true
			})
			wg.Done()
		})

	}
	wg.Wait()
}
