/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/binary"
	"fmt"
	"io"
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
		
		server.MakeTCPServer("Problem 2", meansToAnEndHandler)
	},
}

func init() {
	rootCmd.AddCommand(p2Cmd)
}

func meansToAnEndHandler(conn net.Conn, logger *slog.Logger) {
	defer conn.Close()
	addr := conn.RemoteAddr()
	logger.Info(fmt.Sprintf("New Connection to %s", addr))
	buf := make([]byte, 9)
	clientAssests := make(map[int32]int32)	
	for {
		if _, err := io.ReadFull(conn, buf); err == io.EOF {
			break
		} else if err != nil {
			break
		}
		
		first := int32(binary.BigEndian.Uint32(buf[1:5]))
		second := int32(binary.BigEndian.Uint32(buf[5:]))

		if buf[0] == 'I' {
			clientAssests[first] = second
			logger.Info(fmt.Sprintf("(%d)= %d\n", first, second))
		} else if buf[0] == 'Q' {
			var sum, count, avg int
			for time, val := range clientAssests {
				if first <= time && time <= second {
					sum += int(val)
					count++
				}
			}

			if count > 0 {
				avg = sum / count
			}

			output := make([]byte, 4)
			binary.BigEndian.PutUint32(output, uint32(avg))
			if _, err := conn.Write(output); err != nil {
				logger.Error(fmt.Sprintf("(%s) %s", err, addr))
			} else {
				logger.Info(fmt.Sprintf("query: %v %v ⇒ %v (%v)", first, second, output, addr))
			}
		} else {

			logger.Warn(fmt.Sprintf("Undefined Input for %s to %c", addr, buf[0]))
		}

	}
	logger.Info(fmt.Sprintf("Ended Connection to %s", addr))

}

