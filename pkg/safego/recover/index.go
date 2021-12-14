package recover

import "fmt"

// SafeGo
func SafeGo(f func()) {
	go func() {
		defer func() {
			if panicMessage := recover(); panicMessage != nil {
				fmt.Println(panicMessage)
			}
		}()
		f()
	}()
}
