package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbPath string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dagenie",
		Short: "Dagenie – Purpose-Built for DAGs. Engineered for Speed.",
	}

	rootCmd.PersistentFlags().StringVarP(&dbPath, "db", "d", "./dagdb", "Path to BadgerDB directory")
	connectCmd.Flags().StringVarP(&host, "host", "s", "localhost", "Address of TCP server")
	connectCmd.Flags().StringVarP(&port, "port", "p", "9090", "Port of TCP server")
	serveCmd.Flags().StringVar(&servePort, "port", "9090", "Port to run the TCP server on")
	serveCmd.Flags().StringVar(&dbPath, "db", "", "Path to the database directory")
	serveCmd.MarkFlagRequired("db")
	deleteCmd.Flags().StringVarP(&deleteID, "id", "i", "", "Task ID to delete")
	deleteCmd.Flags().StringVarP(&dagID, "dag", "", "", "DAG ID (required)")
	deleteCmd.MarkFlagRequired("id")
	deleteCmd.MarkFlagRequired("dag")
	insertCmd.Flags().StringVarP(&taskID, "id", "i", "", "Task ID")
	insertCmd.Flags().StringVarP(&taskName, "name", "n", "", "Task name")
	insertCmd.Flags().StringVarP(&payload, "payload", "p", "", "Payload JSON data")
	insertCmd.Flags().StringVarP(&status, "status", "s", "pending", "Task status (e.g., pending, completed)")
	insertCmd.Flags().StringSliceVar(&deps, "deps", []string{}, "Dependencies (comma-separated task IDs)")
	insertCmd.Flags().StringVar(&dagID, "dag", "", "DAG ID (required)")
	insertCmd.MarkFlagRequired("dag")
	insertCmd.MarkFlagRequired("id")
	insertCmd.MarkFlagRequired("name")
	queryCmd.Flags().StringVar(&dqlQuery, "dql", "", "DQL query to execute")
	queryCmd.MarkFlagRequired("dql")
	traverseCmd.Flags().StringVarP(&traverseRoot, "root", "r", "", "Root task ID to start traversal")
	traverseCmd.Flags().StringVarP(&traverseMode, "mode", "m", "dfs", "Traversal mode: dfs or bfs")
	traverseCmd.Flags().StringVar(&dbPath, "db", "", "Path to database (required)")
	traverseCmd.MarkFlagRequired("root")
	traverseCmd.MarkFlagRequired("db")

	// Register commands
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(insertCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(traverseCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(traverseCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
