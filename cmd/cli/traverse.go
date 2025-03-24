package main

import (
	"dagenie/internal/dagdb"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	traverseRoot string
	traverseMode string
)

var traverseCmd = &cobra.Command{
	Use:   "traverse",
	Short: "Traverse DAG from a root task (DFS or BFS)",
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

		// Load tasks into graph
		tasks, err := db.ListAllTasks()
		if err != nil {
			fmt.Println("❌ Failed to load tasks into graph:", err)
			os.Exit(1)
		}
		for _, task := range tasks {
			db.Graph().AddTask(task)
		}

		fmt.Printf("🔍 Traversing DAG from root: %s using %s\n", traverseRoot, traverseMode)

		switch traverseMode {
		case "dfs":
			result := db.Graph().DFS(traverseRoot)
			printTraversal(result)
		case "bfs":
			result := db.Graph().BFS(traverseRoot)
			printTraversal(result)
		default:
			fmt.Println("❌ Invalid traversal mode. Use 'dfs' or 'bfs'.")
			os.Exit(1)
		}
	},
}

func printTraversal(tasks []*dagdb.DAGTask) {
	if len(tasks) == 0 {
		fmt.Println("⚠️ No tasks found during traversal.")
		return
	}
	for _, task := range tasks {
		fmt.Printf("➡️  Task ID=%s Name=%s Status=%s DAGID=%s\n", task.ID, task.Name, task.Status, task.DAGID)
	}
	fmt.Println("✅ Traversal complete.")
}
