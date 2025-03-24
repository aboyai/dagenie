package executor

import (
	"dagenie/internal/dagdb"
	"dagenie/internal/dql/ast"
	"fmt"
)

// ExecuteDelete deletes tasks matching WHERE conditions.
func ExecuteDelete(db *dagdb.DAGDB, deleteAST *ast.DeleteQueryAST) (string, error) {
	if deleteAST.Table != "dag" {
		return "", fmt.Errorf("unsupported table: %s", deleteAST.Table)
	}

	// 1. Load candidate tasks
	var tasks []dagdb.DAGTask
	var err error

	if objectID, ok := deleteAST.Conditions["_id"]; ok {
		tasks, err = db.QueryByObjectID(objectID)
	} else if dagID, ok := deleteAST.Conditions["dagid"]; ok {
		tasks, err = db.ListTasksByDAG(dagID)
	} else {
		tasks, err = db.ListAllTasks()
	}
	if err != nil {
		return "", fmt.Errorf("❌ Task fetch error: %v", err)
	}

	// 2. Filter tasks by WHERE conditions
	var deletedCount int
	for _, task := range tasks {
		match := true
		for field, expected := range deleteAST.Conditions {
			switch field {
			case "dagid":
				if task.DAGID != expected {
					match = false
				}
			case "id":
				if task.ID != expected {
					match = false
				}
			case "status":
				if task.Status != expected {
					match = false
				}
			case "_id":
				if task.ObjectID != expected {
					match = false
				}
			default:
				match = false
			}
			if !match {
				break
			}
		}

		if match {
			err := db.DeleteTask(task.DAGID, task.ID)
			if err != nil {
				fmt.Printf("❌ Failed to delete task: ID=%s, DAGID=%s: %v\n", task.ID, task.DAGID, err)
			} else {
				fmt.Printf("🗑️ Deleted: ID=%s DAGID=%s\n", task.ID, task.DAGID)
				deletedCount++
			}
		}
	}

	if deletedCount == 0 {
		return "❌ No tasks matched for deletion", nil
	}

	return fmt.Sprintf("✅ Deleted %d task(s)", deletedCount), nil
}
