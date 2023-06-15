package safe

import (
	"fmt"
	"log"
	"runtime/debug"
)

// Go 安全go程.
func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(fmt.Sprintf("recover err: %v", err))
				debug.PrintStack()
			}
		}()
		f()
	}()
}

// GoWithField 安全go程且携带参数.
func GoWithField(f func(val interface{}), val interface{}) {
	go func(val interface{}) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(fmt.Sprintf("recover err: %v", err))
				debug.PrintStack()
			}
		}()
		f(val)
	}(val)
}
