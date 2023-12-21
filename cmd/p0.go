/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/rjprice04/protohack/server"
	"github.com/spf13/cobra"
)

// p0Cmd represents the p0 command
var p0Cmd = &cobra.Command{
	Use:   "p0",
	Short: "An echo server",
	Long: `Starts and runs a TCP echo service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("p0 called")
		server.MakeTCPServer(echoServerHandler)
	
	},
}

func init() {
	rootCmd.AddCommand(p0Cmd)
}


func echoServerHandler(conn net.Conn, logger *slog.Logger) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, _ := conn.Read(buffer)

		if n == 0 {
			fmt.Println("Closing Connection")
			break
		}
		conn.Write(buffer[0:n])
	}

}
