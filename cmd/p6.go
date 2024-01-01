/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/rjprice04/protohack/server"
	"github.com/spf13/cobra"
)

// p6Cmd represents the p6 command
var p6Cmd = &cobra.Command{
	Use:   "p6",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("p6 called")
		server.MakeTCPServer("Speed Daemon", speedDeamonHandler)
	},
}

func init() {
	rootCmd.AddCommand(p6Cmd)
}

func speedDeamonHandler(conn net.Conn, logger *slog.Logger) {

}
