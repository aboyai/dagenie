package main

import (
	"dagenie/internal/dagdb"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deleteID string
	dagID    string
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DAG task by ID and DAG",
	Run: func(cmd *cobra.Command, args []string) {
		if dbPath == "" {
			fmt.Println("❌ Please specify the database path using --db")
			return
		}

		db, err := dagdb.OpenDAGDB(dbPath)
		if err != nil {
			fmt.Printf("❌ Error opening DB: %v\n", err)
			return
		}
		defer db.Close()

		err = db.DeleteTask(dagID, deleteID)
		if err != nil {
			fmt.Printf("❌ Delete failed: %v\n", err)
			return
		}
		fmt.Printf("✅ Deleted task: %s from DAG: %s\n", deleteID, dagID)
	},
}
