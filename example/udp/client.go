package main

import (
	"log"
	"net"
)

// client 客户端.
func client() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 6666,
	})
	if err != nil {
		log.Printf("client udp err: %v", err)
		return
	}
	defer socket.Close()

	sendData := []byte("client health")
	if _, err := socket.Write(sendData); err != nil {
		log.Printf("client send udp data err: %v", err)
		return
	}

	data := make([]byte, 2048)
	read, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		log.Printf("client read data err: %v", err)
		return
	}
	log.Printf("client recevie read: %d, addr: %s, data: %s", read, remoteAddr.String(), string(data))
}
