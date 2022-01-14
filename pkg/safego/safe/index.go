package safe

import (
	"log"
)

// Go
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
