/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"

	"github.com/rjprice04/protohack/server"
	"github.com/spf13/cobra"
)

// p1Cmd represents the p1 command
var p1Cmd = &cobra.Command{
	Use:   "p1",
	Short: "Protohackers Problem 2",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		server.MakeTCPServer(primeTimeServerHandler)

		jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)
		myslog.Info("message")
	},
}

func init() {
	rootCmd.AddCommand(p1Cmd)
}


type PrimeTimeOutput struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}


func primeTimeServerHandler(conn net.Conn, logger *slog.Logger) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, _ := conn.Read(buffer)

		if n == 0 {
			logger.Info("Disconnecting", "addr", conn.RemoteAddr())
			break
		}

		if !json.Valid(buffer[0:n]) {
			logger.Error("JSON input was invalid")
			conn.Write([]byte(`{"method":"invalid"}`))
			break
		}
		var input map[string]any
		json.Unmarshal(buffer[0:n], &input)

		method, methodOk := input["method"].(string)
		number, numberOk := input["number"].(float64)

		if !methodOk || !numberOk || method != "isPrime" {
			logger.Error(fmt.Sprintf("JSON input was invalid (method ok: %t, number ok: %t)", methodOk, numberOk))
			conn.Write([]byte(`{"method":"invalid"}`))
			break
		}

		var output PrimeTimeOutput
		output.Method = "isPrime"
		// Check if a float is actual an int i.e. 4.0 can be 4 but 4.5 is not a valid int
		if number == float64(int(number)) {
			output.Prime = big.NewInt(int64(number)).ProbablyPrime(0)
		} else {
			// floats are always false
			output.Prime = false
		}
		outBuf, _ := json.Marshal(output)
		
		outBuf = append(outBuf, '\r')
		outBuf = append(outBuf, '\n')
		conn.Write(outBuf)

	}

}
