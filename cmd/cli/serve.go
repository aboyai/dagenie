package main

import (
	"dagenie/internal/dagdb"
	"dagenie/internal/tcp"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var servePort string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Dagenie TCP server",
	Run: func(cmd *cobra.Command, args []string) {
		if dbPath == "" {
			fmt.Println("❌ Please provide a database path using --db flag")
			os.Exit(1)
		}

		db, err := dagdb.OpenDAGDB(dbPath)
		if err != nil {
			fmt.Printf("❌ Error opening DB at '%s': %v\n", dbPath, err)
			os.Exit(1)
		}
		defer db.Close()

		address := ":" + servePort
		fmt.Printf("🚀 Starting Dagenie server on port %s using DB: %s\n", servePort, dbPath)

		err = tcp.StartTCPServer(db, address)
		if err != nil {
			fmt.Println("❌ TCP Server error:", err)
			os.Exit(1)
		}
	},
}

