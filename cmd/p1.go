/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net"

	"github.com/rjprice04/protohack/server"
	"github.com/spf13/cobra"
)

// p1Cmd represents the p1 command
var p1Cmd = &cobra.Command{
	Use:   "p1",
	Short: "Protohackers Problem 1",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		
		server.MakeTCPServer("Problem 1", primeTimeServerHandler)

	},
}

func init() {
	rootCmd.AddCommand(p1Cmd)
}
type Request struct {
	Method string `json:"method"`
	Number *float64 `json:"number"`
}

type PrimeTimeOutput struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}


func primeTimeServerHandler(conn net.Conn, logger *slog.Logger) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		buffer := scanner.Bytes()


		logger.Info(fmt.Sprintf("body: %s", string(buffer)))

		var req Request
		err := json.Unmarshal(buffer, &req)

		if err !=nil || req.Method != "isPrime" || req.Number == nil {
			conn.Write([]byte(`{"method":"invalid"}`))
			break
		}

		var output PrimeTimeOutput
		output.Method = "isPrime"
		// Check if a float is actual an int i.e. 4.0 can be 4 but 4.5 is not a valid int
		if *req.Number == float64(int(*req.Number)) {
			output.Prime = big.NewInt(int64(*req.Number)).ProbablyPrime(0)
		} else {
			// floats are always false
			output.Prime = false
		}
		outBuf, _ := json.Marshal(output)
		
		outBuf = append(outBuf, '\n')
		conn.Write(outBuf)

	}

}
