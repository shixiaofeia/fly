package recover

import (
	"log"
)

// SafeGo
func SafeGo(f func()) {
	go func() {
		defer func() {
			if msg := recover(); msg != nil {
				log.Println(msg)
			}
		}()
		f()
	}()
}
