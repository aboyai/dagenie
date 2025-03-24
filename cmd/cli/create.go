package main

import (
	"dagenie/internal/dagdb"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create/init DAG database",
	Run: func(cmd *cobra.Command, args []string) {
		if dbPath == "" {
			fmt.Println("❌ Please specify the database path using --db")
			return
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
				fmt.Println("❌ Error creating DB directory:", err)
				return
			}
		}

		db, err := dagdb.OpenDAGDB(dbPath)
		if err != nil {
			fmt.Println("❌ Failed to create DB:", err)
			return
		}
		defer db.Close()

		fmt.Println("✅ DAG database initialized at:", dbPath)
	},
}
