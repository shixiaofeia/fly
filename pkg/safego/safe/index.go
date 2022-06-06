package safe

import (
	"log"
)

// Go 安全go程.
func Go(f func()) {
	go func() {
		defer func() {
			if msg := recover(); msg != nil {
				log.Println(msg)
			}
		}()
		f()
	}()
}

// GoWithField 安全go程且携带参数
func GoWithField(f func(val interface{}), val interface{}) {
	go func(val interface{}) {
		defer func() {
			if msg := recover(); msg != nil {
				log.Println(msg)
			}
		}()
		f(val)
	}(val)
}
