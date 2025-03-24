package main

import (
	"encoding/json"
	"fmt"

	"dagenie/internal/dagdb"
	"dagenie/utils"

	"github.com/spf13/cobra"
)

var (
	taskID   string
	taskName string
	payload  string
	status   string
	deps     []string
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a new DAG task",
	Run: func(cmd *cobra.Command, args []string) {
		if dbPath == "" {
			fmt.Println("❌ Please specify the database path using --db")
			return
		}

		db, err := dagdb.OpenDAGDB(dbPath)
		if err != nil {
			fmt.Println("❌ Error opening DB:", err)
			return
		}
		defer db.Close()

		// Validate Payload is valid JSON
		var js json.RawMessage
		if err := json.Unmarshal([]byte(payload), &js); err != nil {
			fmt.Println("❌ Invalid JSON payload")
			return
		}

		// Create task struct
		task := dagdb.DAGTask{
			ObjectID:     utils.GenerateObjectID(),
			ID:           taskID,
			Name:         taskName,
			Payload:      payload,
			Status:       status,
			DAGID:        dagID,
			Dependencies: deps,
		}

		if err := db.SaveTaskWithCycleCheck(task); err != nil {
			fmt.Printf("❌ Insert failed: %v\n", err)
			return
		}
		fmt.Printf("✅ Task inserted: %s\n", task.ID)
	},
}
