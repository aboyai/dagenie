package ast

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// LogicalNode is the interface all logical expression nodes implement.
type LogicalNode interface {
	Evaluate(task interface{}) bool
}

// ---------------- ConditionNode ------------------

// ConditionNode is a leaf node in the logical tree.
type ConditionNode struct {
	Field    string
	Operator string // e.g., "=", "!=", ">", "<", "<=", ">="
	Value    string
}

func (c *ConditionNode) Evaluate(task interface{}) bool {
	v := reflect.ValueOf(task)

	// Get the field value via reflection
	fieldVal := v.FieldByNameFunc(func(name string) bool {
		return strings.EqualFold(name, c.Field)
	})

	if !fieldVal.IsValid() {
		fmt.Printf("❌ Field '%s' not found in task\n", c.Field)
		return false
	}

	var taskValStr string
	switch fieldVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		taskValStr = strconv.FormatInt(fieldVal.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		taskValStr = strconv.FormatUint(fieldVal.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		taskValStr = fmt.Sprintf("%.2f", fieldVal.Float())
	case reflect.String:
		taskValStr = fieldVal.String()
	default:
		return false
	}

	// Perform comparison
	switch c.Operator {
	case "=":
		return taskValStr == c.Value
	case "!=":
		return taskValStr != c.Value
	case ">":
		return taskValStr > c.Value
	case "<":
		return taskValStr < c.Value
	case ">=":
		return taskValStr >= c.Value
	case "<=":
		return taskValStr <= c.Value
	default:
		return false
	}
}

// ---------------- AndNode ------------------

type AndNode struct {
	Left, Right LogicalNode
}

func (a *AndNode) Evaluate(task interface{}) bool {
	return a.Left.Evaluate(task) && a.Right.Evaluate(task)
}

// ---------------- OrNode ------------------

type OrNode struct {
	Left, Right LogicalNode
}

func (o *OrNode) Evaluate(task interface{}) bool {
	return o.Left.Evaluate(task) || o.Right.Evaluate(task)
}

// ---------------- NotNode ------------------

type NotNode struct {
	Expr LogicalNode
}

func (n *NotNode) Evaluate(task interface{}) bool {
	return !n.Expr.Evaluate(task)
}
