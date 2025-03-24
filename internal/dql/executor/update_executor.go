package executor

import (
	"dagenie/internal/dagdb"
	"dagenie/internal/dql/ast"
	"fmt"
	"strconv"
	"strings"
)

func ExecuteUpdate(db *dagdb.DAGDB, updateAST *ast.UpdateQueryAST) (string, error) {
	if updateAST.Table != "dag" {
		return "", fmt.Errorf("❌ Unsupported table: %s", updateAST.Table)
	}

	// Load all tasks into memory
	tasks, err := db.ListAllTasks()
	if err != nil {
		return "", fmt.Errorf("❌ Task load error: %v", err)
	}

	updatedCount := 0

	for _, task := range tasks {
		// Check WHERE condition
		if !matchesCondition(task, updateAST.Where) {
			continue
		}

		oldTask := task // For key comparison
		updated := false

		// Apply SET clause
		for field, value := range updateAST.SetFields {
			switch strings.ToLower(field) {
			case "id":
				if task.ID != value {
					task.ID = value
					updated = true
				}
			case "dagid":
				if task.DAGID != value {
					task.DAGID = value
					updated = true
				}
			case "name":
				if task.Name != value {
					task.Name = value
					updated = true
				}
			case "status":
				if task.Status != value {
					task.Status = value
					updated = true
				}
			case "payload":
				if task.Payload != value {
					task.Payload = value
					updated = true
				}
			case "duration":
				dur, err := strconv.Atoi(value)
				if err != nil {
					return "", fmt.Errorf("❌ Invalid duration value: %v", err)
				}
				if task.Duration != dur {
					task.Duration = dur
					updated = true
				}
			case "retries":
				ret, err := strconv.Atoi(value)
				if err != nil {
					return "", fmt.Errorf("❌ Invalid retries value: %v", err)
				}
				if task.Retries != ret {
					task.Retries = ret
					updated = true
				}
			default:
				return "", fmt.Errorf("❌ Unknown field: %s", field)
			}
		}

		if updated {
			if task.ID != oldTask.ID || task.DAGID != oldTask.DAGID {
				// Key changed → migrate key
				err := db.UpdateTaskWithKeyChange(oldTask, task)
				if err != nil {
					return "", fmt.Errorf("❌ Key migration error: %v", err)
				}
			} else {
				// Update task in DB and graph
				err := db.SaveTask(task)
				if err != nil {
					return "", fmt.Errorf("❌ Save error: %v", err)
				}
			}
			// Always update the graph structure
			db.UpdateGraphTask(&task)
			updatedCount++
		}
	}

	if updatedCount == 0 {
		return "❌ No matching tasks found", nil
	}

	return fmt.Sprintf("✅ Updated %d task(s)", updatedCount), nil
}

// matchesCondition checks WHERE clause against task fields
func matchesCondition(task dagdb.DAGTask, where map[string]string) bool {
	for field, expected := range where {
		switch strings.ToLower(field) {
		case "dagid":
			if strings.ToLower(task.DAGID) != strings.ToLower(expected) {
				return false
			}
		case "id":
			if strings.ToLower(task.ID) != strings.ToLower(expected) {
				return false
			}
		case "name":
			if strings.ToLower(task.Name) != strings.ToLower(expected) {
				return false
			}
		case "status":
			if strings.ToLower(task.Status) != strings.ToLower(expected) {
				return false
			}
		case "payload":
			if strings.ToLower(task.Payload) != strings.ToLower(expected) {
				return false
			}
		case "_id":
			if strings.ToLower(task.ObjectID) != strings.ToLower(expected) {
				return false
			}
		default:
			// Unknown fields are ignored here
			return false
		}
	}
	return true
}
