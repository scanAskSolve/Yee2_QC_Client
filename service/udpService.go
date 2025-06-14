package main

import (
	"log"
	"net"
	"time"
)

func udpService() []string {

	foundIps := make(map[string]struct{})

	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Println("Listening for UDP broadcasts on port 9999...")

	buffer := make([]byte, 1024)
	timeout := time.After(10 * time.Second)

	for {
		select {
		case <-timeout:
			var result []string
			for ip := range foundIps {
				result = append(result, ip)
			}
			return result
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Println("Error reading from UDP:", err)
				continue
			}
			if string(buffer[:n]) == "DISCOVER_DEVICE" {
				log.Printf("Received message from %s: %s\n", remoteAddr.String(), string(buffer[:n]))
				log.Println(remoteAddr.IP.String())
				foundIps[remoteAddr.IP.String()] = struct{}{}
			}
		}
	}

}
