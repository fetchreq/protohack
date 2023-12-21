package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
)

type fn func(conn net.Conn, logger *slog.Logger)

func MakeTCPServer(handler fn) {

	// Set up logger	
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog := slog.New(jsonHandler)
	myslog.Info("message")

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		myslog.Info("New Connection", "addr", conn.RemoteAddr())
		go handler(conn, myslog)
	}
}
