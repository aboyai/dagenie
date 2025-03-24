package dql

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dagenie/internal/dagdb"
	"dagenie/internal/dql/executor"
	"dagenie/internal/dql/parser"
)

var openDBs = make(map[string]*dagdb.DAGDB)

// ---------------------- Dispatch Executor ----------------------

func ExecuteDQLWithContext(db *dagdb.DAGDB, query string) (string, *dagdb.DAGDB, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", nil, fmt.Errorf("empty query")
	}
	lower := strings.ToLower(query)

	switch {
	// CREATE DATABASE
	case strings.HasPrefix(lower, "create database"):
		dbName := strings.TrimSpace(strings.TrimSuffix(query[len("create database"):], ";"))
		if dbName == "" {
			return "", nil, fmt.Errorf("❌ Invalid database name")
		}
		dbPath := filepath.Join("./data", dbName)
		if _, err := os.Stat(dbPath); err == nil {
			return "", nil, fmt.Errorf("❌ Database exists")
		}

		if err := os.MkdirAll(dbPath, 0755); err != nil {
			return "", nil, fmt.Errorf("❌ Create DB failed: %v", err)
		}
		return fmt.Sprintf("✅ Database '%s' created", dbName), nil, nil

		// USE DB
		// USE DB
	case strings.HasPrefix(lower, "use"):
		dbName := strings.TrimSpace(query[4:]) // Use original case for name
		if dbName == "" {
			return "", nil, fmt.Errorf("❌ No database specified")
		}
		dbPath := "./data/" + dbName

		// Check if directory exists
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return "", nil, fmt.Errorf("❌ Database '%s' does not exist", dbName)
		}

		// Check if DB is already open
		if existingDB, ok := openDBs[dbName]; ok {
			return fmt.Sprintf("✅ Using database '%s'", dbName), existingDB, nil
		}

		// Open new DB and cache it
		newDB, err := dagdb.OpenDAGDB(dbPath)
		if err != nil {
			return "", nil, fmt.Errorf("❌ Failed to open DB '%s': %v", dbName, err)
		}
		openDBs[dbName] = newDB // Cache it

		return fmt.Sprintf("✅ Using database '%s'", dbName), newDB, nil

	// SHOW DATABASES
	case strings.HasPrefix(lower, "show databases"):
		entries, err := os.ReadDir("./data")
		if err != nil {
			return "", nil, fmt.Errorf("❌ Failed to read DBs: %v", err)
		}
		var dbs []string
		for _, entry := range entries {
			fmt.Println(entry)

			if entry.IsDir() {
				dbs = append(dbs, entry.Name())
			}
		}

		return fmt.Sprintf("%s %s", strings.Join(dbs, "\n"), "\n✅ Done"), nil, nil

	// DROP DATABASE
	case strings.HasPrefix(lower, "drop database"):
		dbName := strings.TrimSpace(strings.TrimSuffix(query[len("drop database"):], ";"))
		dbPath := filepath.Join("./data", dbName)
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return "", nil, fmt.Errorf("❌ Database '%s' not found", dbName)
		}
		err := os.RemoveAll(dbPath)
		if err != nil {
			return "", nil, fmt.Errorf("❌ Delete failed: %v", err)
		}
		return fmt.Sprintf("🗑️ Database '%s' deleted", dbName), nil, nil

	// Pass to existing DQL (SELECT, INSERT, etc.)
	default:
		result, err := ExecuteDQL(db, query)
		fmt.Println(result)
		return result, nil, err
	}
}

// ExecuteDQL dispatches raw query to parser → executor
func ExecuteDQL(globalDB *dagdb.DAGDB, queryLine string) (string, error) {
	queryLine = strings.TrimSpace(queryLine)
	if queryLine == "" {
		return "", fmt.Errorf("empty query")
	}

	lowerQuery := strings.ToLower(queryLine)

	switch {
	case strings.HasPrefix(lowerQuery, "select"):
		astSelect, err := parser.ParseSelectToAST(queryLine)
		if err != nil {
			return "", fmt.Errorf("❌ SELECT Parse Error: %v", err)
		}
		result, err := executor.ExecuteSelect(globalDB, astSelect)
		if err != nil {
			return "", fmt.Errorf("❌ SELECT Execution Error: %v", err)
		}
		return result, nil

	case strings.HasPrefix(lowerQuery, "insert"):
		insertAST, err := parser.ParseInsertToAST(queryLine)
		if err != nil {
			return "", fmt.Errorf("❌ INSERT Parse Error: %v", err)
		}
		result, err := executor.ExecuteInsert(globalDB, insertAST)
		if err != nil {
			return "", fmt.Errorf("❌ INSERT Execution Error: %v", err)
		}
		return result, nil

	case strings.HasPrefix(lowerQuery, "update"):
		updateAST, err := parser.ParseUpdateToAST(queryLine)
		if err != nil {
			return "", fmt.Errorf("UPDATE Parse Error: %v", err)
		}
		result, err := executor.ExecuteUpdate(globalDB, updateAST)
		if err != nil {
			return "", fmt.Errorf("UPDATE Execution Error: %v", err)
		}
		return result, nil

	case strings.HasPrefix(lowerQuery, "delete"):
		deleteAST, err := parser.ParseDeleteToAST(queryLine)
		if err != nil {
			return "", fmt.Errorf("❌ DELETE Parse Error: %v", err)
		}
		result, err := executor.ExecuteDelete(globalDB, deleteAST)
		if err != nil {
			return "", fmt.Errorf("❌ DELETE Execution Error: %v", err)
		}
		return result, nil

	default:
		return "", fmt.Errorf("❌ Unsupported query type: %s", strings.Split(queryLine, " ")[0])
	}
}
