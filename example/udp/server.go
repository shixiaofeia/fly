package main

import (
	"log"
	"net"
)

// server 服务端
func server() {
	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 6666,
	})
	if err != nil {
		log.Printf("server listen udp err: %v", err)
		return
	}
	defer socket.Close()

	for {
		data := make([]byte, 2048)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			log.Printf("server read data err: %v", err)
			continue
		}
		log.Printf("server recevie read: %d, addr: %s, data: %s", read, remoteAddr.String(), string(data))

		sendData := []byte(`server health`)
		if _, err := socket.WriteToUDP(sendData, remoteAddr); err != nil {
			log.Printf("server write udp data err: %v", err)
		}
	}
}
