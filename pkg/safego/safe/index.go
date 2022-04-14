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
