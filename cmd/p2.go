/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"

	"github.com/rjprice04/protohack/server"
	"github.com/spf13/cobra"
)

// p2Cmd represents the p2 command
var p2Cmd = &cobra.Command{
	Use:   "p2",
	Short: "Means to an End",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("p2 called")
		
		//meansToAnEnd();
		server.MakeTCPServer("Problem 2", primeTimeServerHandler)
	},
}

func init() {
	rootCmd.AddCommand(p2Cmd)
}

func meansToAnEndHandler(conn net.Conn, logger *slog.Logger) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	addr := conn.RemoteAddr()
	logger.Info(fmt.Sprintf("New Connection to %s", addr))
	for scanner.Scan() {
		buffer := scanner.Bytes();

		if len(buffer) != 9 {
			continue
		}

		mode := buffer[0]
		if mode == 'I' {
			// Insert
			timestamp := binary.BigEndian.Uint32(buffer[1:5])
			value := binary.BigEndian.Uint32(buffer[5:])

			fmt.Printf("Insert at %d with value %d", timestamp, value)
		} else if mode == 'Q' {
			minTime := binary.BigEndian.Uint32(buffer[1:5])
			maxTime := binary.BigEndian.Uint32(buffer[5:])

			fmt.Printf("Query between %d and %d", minTime, maxTime)

			//Query 
		} else {
			// Undefined
			continue
		}

	}
	logger.Info(fmt.Sprintf("Ended Connection to %s", addr))

}

func meansToAnEnd() {
	//buffer := []byte{0x49,0x00,0x00,0x30,0x39,0x00,0x00,0x00,0x65}
	//buffer := []byte{0x51,0x00,0x00,0x03,0xe8,0x00,0x01,0x86,0xa0}
	//
	// if len(buffer) != 9 {
	// 	fmt.Println("Not Enough bytes")
	// 	return
	// }
	// //mode := buffer[0]
	// part1 := buffer[1:5]
	// part2 := buffer[5:]
	// if part1[0] == 1 {
	// 	fmt.Println("")
	// }
	//
	// if part2[0] == 1 {
	//
	// }

	// if mode == 'I' {
	// 	// Insert
	// 	timestamp := binary.BigEndian.Uint32(buffer[1:5])
	// 	value := binary.BigEndian.Uint32(buffer[5:])
	//
	// 	fmt.Printf("Insert at %d with value %d", timestamp, value)
	// } else if mode == 'Q' {
	// 	// Insert
	// 	minTime := binary.BigEndian.Uint32(buffer[1:5])
	// 	maxTime := binary.BigEndian.Uint32(buffer[5:])
	// 	
	// 	fmt.Printf("Query between %d and %d", minTime, maxTime)
	//
	// 	//Query 
	// } else {
	// 	fmt.Printf("Invalid First Byte %c", mode)
	// 	// Undefined
	// }
	


}
