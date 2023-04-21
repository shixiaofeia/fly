package main

func main() {
	go server()

	var ch chan struct{}
	client()
	<-ch
}
