package parser

import (
	"dagenie/internal/dql/ast"
	"fmt"
	"regexp"
	"strings"
)

// ParseDeleteToAST parses a DELETE query into DeleteQueryAST.
func ParseDeleteToAST(query string) (*ast.DeleteQueryAST, error) {
	query = strings.TrimSpace(query)
	lower := strings.ToLower(query)

	if !strings.HasPrefix(lower, "delete") {
		return nil, fmt.Errorf("❌ Not a DELETE query")
	}

	fromIdx := strings.Index(lower, "from")
	if fromIdx == -1 {
		return nil, fmt.Errorf("❌ Missing FROM clause")
	}

	// Extract FROM table and optional WHERE
	fromRest := strings.TrimSpace(query[fromIdx+4:])
	if fromRest == "" {
		return nil, fmt.Errorf("❌ Missing table name after FROM")
	}

	whereIdx := strings.Index(strings.ToLower(fromRest), "where")
	var tableName, wherePart string

	if whereIdx != -1 {
		tableName = strings.TrimSpace(fromRest[:whereIdx])
		wherePart = strings.TrimSpace(fromRest[whereIdx+5:])
	} else {
		tableName = fromRest
	}

	if tableName == "" {
		return nil, fmt.Errorf("❌ Missing table name after FROM")
	}

	// Parse WHERE conditions
	conditions := map[string]string{}
	if wherePart != "" {
		andRegex := regexp.MustCompile(`(?i)\s+and\s+`)
		whereClauses := andRegex.Split(wherePart, -1)
		for _, clause := range whereClauses {
			kv := strings.SplitN(clause, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("❌ Invalid WHERE clause: %s", clause)
			}
			key := strings.ToLower(strings.TrimSpace(kv[0]))
			val := strings.Trim(strings.TrimSpace(kv[1]), `"'`)
			conditions[key] = val
		}
	}

	return &ast.DeleteQueryAST{
		Table:      strings.ToLower(tableName),
		Conditions: conditions,
	}, nil
}
