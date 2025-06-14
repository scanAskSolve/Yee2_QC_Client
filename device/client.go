package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		message := []byte("DISCOVER_DEVICE")
		_, err := conn.Write(message)
		if err != nil {
			fmt.Println("Error sending:", err)
		} else {
			fmt.Println("Broadcast sent")
		}
		time.Sleep(5 * time.Second)
	}
}
