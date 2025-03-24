package main

import (
	"fmt"
	"os"

	"dagenie/internal/dagdb"
	"dagenie/internal/dql"

	"github.com/spf13/cobra"
)

var dqlQuery string

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Execute a Dagenie DQL query",
	Run: func(cmd *cobra.Command, args []string) {
		if dbPath == "" {
			fmt.Println("❌ Please provide a database path using --db flag")
			os.Exit(1)
		}

		db, err := dagdb.OpenDAGDB(dbPath)
		if err != nil {
			fmt.Printf("❌ Failed to open DB: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		if dqlQuery == "" {
			fmt.Println("❌ Please provide a DQL query using --dql flag")
			os.Exit(1)
		}

		result, err := dql.ExecuteDQL(db, dqlQuery)
		if err != nil {
			fmt.Printf("❌ Execution Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(result)
	},
}

