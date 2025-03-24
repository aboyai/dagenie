package parser

import (
	"dagenie/internal/dql/ast"
	"fmt"
	"strings"
)

type UpdateQuery struct {
	Table      string
	SetFields  map[string]string
	Conditions map[string]string
}

func ParseUpdateToAST(query string) (*ast.UpdateQueryAST, error) {
	// Trim and lowercase
	query = strings.TrimSpace(query)
	lower := strings.ToLower(query)

	if !strings.HasPrefix(lower, "update") {
		return nil, fmt.Errorf("❌ Not an UPDATE query")
	}

	setIdx := strings.Index(lower, "set")
	whereIdx := strings.Index(lower, "where")

	if setIdx == -1 {
		return nil, fmt.Errorf("❌ Missing SET clause")
	}

	table := strings.TrimSpace(query[6:setIdx])
	if table == "" {
		return nil, fmt.Errorf("❌ Missing table name")
	}

	var setPart, wherePart string
	if whereIdx != -1 {
		setPart = query[setIdx+3 : whereIdx]
		wherePart = query[whereIdx+5:]
	} else {
		setPart = query[setIdx+3:]
	}

	// Parse SET
	setFields := make(map[string]string)
	assignments := strings.Split(setPart, ",")
	for _, assign := range assignments {
		kv := strings.SplitN(assign, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("❌ Invalid SET clause: %s", assign)
		}
		field := strings.ToLower(strings.TrimSpace(kv[0]))
		value := strings.Trim(strings.TrimSpace(kv[1]), `"'`)
		setFields[field] = value
	}

	// Parse WHERE
	where := make(map[string]string)
	if wherePart != "" {
		conditions := strings.Split(wherePart, "and")
		for _, cond := range conditions {
			kv := strings.SplitN(cond, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("❌ Invalid WHERE clause: %s", cond)
			}
			field := strings.ToLower(strings.TrimSpace(kv[0]))
			value := strings.Trim(strings.TrimSpace(kv[1]), `"'`)
			where[field] = value
		}
	}

	return &ast.UpdateQueryAST{
		Table:     strings.ToLower(table),
		SetFields: setFields,
		Where:     where, // ✅ Now exists
	}, nil
}
