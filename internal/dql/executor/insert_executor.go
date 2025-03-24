package executor

import (
	"dagenie/internal/dagdb"
	"dagenie/internal/dql/ast"
	"dagenie/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ExecuteInsert(db *dagdb.DAGDB, insertAST *ast.InsertQueryAST) (string, error) {
	if insertAST.Table != "dag" {
		return "", fmt.Errorf("❌ Unsupported table: %s", insertAST.Table)
	}

	// Map columns to values (lowercased keys)
	data := make(map[string]string)
	for i, col := range insertAST.Columns {
		data[strings.ToLower(col)] = insertAST.Values[i]
	}

	// Validate required fields
	requiredFields := []string{"id", "name", "status", "payload", "dependencies", "dagid", "duration", "retries"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			return "", fmt.Errorf("❌ Missing required field: %s", field)
		}
	}

	// Validate 'dagid' - no spaces
	dagid := data["dagid"]
	if strings.Contains(dagid, " ") {
		return "", fmt.Errorf("❌ Invalid value for dagid: cannot contain spaces")
	}

	// Validate 'name' - no spaces
	name := data["name"]
	if strings.Contains(name, " ") {
		return "", fmt.Errorf("❌ Invalid value for name: cannot contain spaces")
	}

	// Convert duration and retries to int
	durationInt, err := strconv.Atoi(data["duration"])
	if err != nil {
		return "", fmt.Errorf("❌ Invalid duration value: %v", err)
	}

	retriesInt, err := strconv.Atoi(data["retries"])
	if err != nil {
		return "", fmt.Errorf("❌ Invalid retries value: %v", err)
	}

	// Parse dependencies - must be a JSON array string
	var dependencies []string
	err = json.Unmarshal([]byte(data["dependencies"]), &dependencies)
	if err != nil {
		return "", fmt.Errorf("❌ Invalid dependencies format: %v", err)
	}

	// Create task object
	task := dagdb.DAGTask{
		ObjectID:     utils.GenerateObjectID(),
		DAGID:        dagid,
		ID:           data["id"],
		Name:         name,
		Payload:      data["payload"],
		Status:       data["status"],
		Duration:     durationInt,
		Retries:      retriesInt,
		Dependencies: dependencies,
	}

	// Save to database
	err = db.SaveTask(task)
	if err != nil {
		return "", fmt.Errorf("❌ Insert Failed: %v", err)
	}
	// ✅ Add to in-memory graph
	db.Graph().AddTask(task)

	fmt.Println("I8")
	fmt.Println("📊 Current Graph Size:", len(db.Graph().AllTasks()))

	return fmt.Sprintf("✅ Inserted task ID=%s", task.ID), nil
}
