package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := map[string]int{"data": 1}
	json.NewEncoder(w).Encode(resp)
}

func main() {

	logger := &lumberjack.Logger{
		Filename:   "./logs/client.log", // log 檔案路徑
		MaxSize:    10,                  // 每個 log 檔最大 MB
		MaxBackups: 200,                 // 最多保留幾個舊 log
		MaxAge:     60,                  // 最多保留幾天
		Compress:   false,               // 是否壓縮舊 log
		LocalTime:  true,                // 使用本地時間分割
	}
	multiWriter := io.MultiWriter(os.Stdout, logger)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("service started")

	go func() {
		http.HandleFunc("/health", healthHandler)
		log.Println("Http service listening on 8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}

	}()

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
