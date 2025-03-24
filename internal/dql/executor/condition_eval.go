package executor

import (
	"dagenie/internal/dagdb"
	"dagenie/internal/dql/ast"
	"fmt"
	"strconv"
	"strings"
)

// Entry point for evaluating logical tree on a task
func evaluateConditionTree(task dagdb.DAGTask, node ast.LogicalNode) bool {
	if node == nil {
		return true // No condition
	}
	return node.Evaluate(task)
}

func getField(task dagdb.DAGTask, field string) string {
	switch strings.ToLower(field) {
	case "id":
		return task.ID
	case "name":
		return task.Name
	case "status":
		return task.Status
	case "dagid":
		return task.DAGID
	case "_id":
		return task.ObjectID
	case "payload":
		return task.Payload
	case "duration":
		return fmt.Sprintf("%d", task.Duration)
	case "retries":
		return fmt.Sprintf("%d", task.Retries)
	default:
		return ""
	}
}

func evaluateCondition(task dagdb.DAGTask, cond *ast.ConditionNode) bool {
	field := strings.ToLower(cond.Field)
	value := cond.Value

	switch field {
	case "dagid":
		return task.DAGID == value
	case "id":
		return task.ID == value
	case "name":
		return task.Name == value
	case "status":
		return task.Status == value
	case "payload":
		return task.Payload == value
	case "duration":
		dur, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		return task.Duration == dur
	case "retries":
		retries, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		return task.Retries == retries
	case "_id":
		return task.ObjectID == value
	default:
		return false
	}
}

func evaluateSingleCondition(task dagdb.DAGTask, cond *ast.ConditionNode) bool {
	val := strings.ToLower(cond.Value)

	switch cond.Field {
	case "dagid":
		return strings.ToLower(task.DAGID) == val
	case "id":
		return strings.ToLower(task.ID) == val
	case "name":
		return strings.ToLower(task.Name) == val
	case "status":
		return strings.ToLower(task.Status) == val
	case "payload":
		return strings.ToLower(task.Payload) == val
	case "duration":
		intVal, _ := strconv.Atoi(val)
		return task.Duration == intVal
	case "retries":
		intVal, _ := strconv.Atoi(val)
		return task.Retries == intVal
	case "_id":
		return strings.ToLower(task.ObjectID) == val
	default:
		return false
	}
}
