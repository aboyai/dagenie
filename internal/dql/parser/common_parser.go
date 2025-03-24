package parser

import (
	"fmt"
	"strings"
)

// parseFromWhere extracts the table name and WHERE conditions from the FROM clause.
func parseFromWhere(fromRest string) (string, map[string]string, error) {
	fromRestLower := strings.ToLower(fromRest)

	// Find WHERE clause if present
	whereIdx := strings.Index(fromRestLower, "where")

	tablePart := fromRest
	wherePart := ""

	if whereIdx != -1 {
		tablePart = strings.TrimSpace(fromRest[:whereIdx])
		wherePart = strings.TrimSpace(fromRest[whereIdx+5:])
	}

	// Extract clean table name (no trailing group/order/limit)
	tableTokens := strings.Fields(tablePart)
	if len(tableTokens) == 0 {
		return "", nil, fmt.Errorf("❌ Missing table name")
	}
	tableName := tableTokens[0]

	// Parse WHERE conditions (basic field='value' AND ...)
	conditions := make(map[string]string)
	if wherePart != "" {
		clauses := strings.Split(wherePart, "and")
		for _, clause := range clauses {
			parts := strings.SplitN(clause, "=", 2)
			if len(parts) != 2 {
				return "", nil, fmt.Errorf("❌ Invalid WHERE clause: %s", clause)
			}
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove surrounding quotes
			value = strings.Trim(value, "'\"")
			conditions[strings.ToLower(field)] = value
		}
	}

	return strings.ToLower(tableName), conditions, nil
}
