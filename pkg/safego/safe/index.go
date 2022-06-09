package safe

import (
	"fmt"
	"log"
	"runtime"
)

// Go 安全go程.
func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				_, file, line, _ := runtime.Caller(3)
				log.Println(fmt.Sprintf("recover err: %v, line: %s", err, fmt.Sprintf("%s:%d", file, line)))
			}
		}()
		f()
	}()
}

// GoWithField 安全go程且携带参数
func GoWithField(f func(val interface{}), val interface{}) {
	go func(val interface{}) {
		defer func() {
			if err := recover(); err != nil {
				_, file, line, _ := runtime.Caller(3)
				log.Println(fmt.Sprintf("recover err: %v, line: %s", err, fmt.Sprintf("%s:%d", file, line)))
			}
		}()
		f(val)
	}(val)
}
