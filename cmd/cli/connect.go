package main

import (
	"dagenie/internal/tcp"

	"github.com/spf13/cobra"
)

var host string
var port string

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to Dagenie TCP Client",
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr := host + ":" + port
		tcp.StartTCPClient(serverAddr)
	},
}
