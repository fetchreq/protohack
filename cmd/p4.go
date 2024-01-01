/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// p4Cmd represents the p4 command
var p4Cmd = &cobra.Command{
	Use:   "p4",
	Short: "Unusual Database Program",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Set up logger	
		jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)
		myslog.Info("Starting Unusual Database Program")

		packetConn, err := net.ListenPacket("udp", "fly-global-services:10000")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer packetConn.Close()
		buffer := make([]byte, 1000)
		database := make(map[string]string)
		database["version"] = "KV 1.0"
		for {
			n, addr, err :=  packetConn.ReadFrom(buffer)
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}

			message := string(buffer[0:n-1])
			if strings.Contains(message, "=") {

				messageParts := strings.SplitN(message, "=", 2)

				if messageParts[0] != "version" {
					fmt.Printf("key = (%s) val = (%s)\n", messageParts[0], messageParts[1])
					database[messageParts[0]] = messageParts[1]
				} else {
					fmt.Println("version update not allowed")
				}
			} else {
				fmt.Printf("look up for key = (%s)\n", message)
				value, ok := database[message]
				if ok {
					fmt.Printf("key = %s val = %s\n", message, value)

					output := []byte(fmt.Sprintf("%s=%s", message, value))
					

					packetConn.WriteTo(output, addr)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(p4Cmd)
}
