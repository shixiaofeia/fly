package gopool

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestNewGoPool(t *testing.T) {
	defaultNum := runtime.NumGoroutine()
	p := NewGoPool(2)

	for i := 0; i < 100; i++ {
		p.Add()
		go func() {
			defer p.Done()
			time.Sleep(1 * time.Second)
			fmt.Println("go routine num: ", runtime.NumGoroutine()-defaultNum)
		}()
	}
	p.Wait()
	fmt.Println("go routine num: ", runtime.NumGoroutine())
}
