package parser

import (
	"dagenie/internal/dql/ast"
	"fmt"
	"regexp"
	"strings"
)

var orderAggRegex = regexp.MustCompile(`(?i)(sum|avg|max|min|count)\s*\(\s*([a-zA-Z0-9_*]+)\s*\)`)
var aggregateRegex = regexp.MustCompile(`(?i)(sum|avg|max|min|count)\s*\(\s*([a-zA-Z0-9_*]+)\s*\)`)

func ParseSelectToAST(query string) (*ast.SelectQueryAST, error) {
	query = strings.TrimSpace(query)
	lowerQuery := strings.ToLower(query)

	if !strings.HasPrefix(lowerQuery, "select") {
		return nil, fmt.Errorf("❌ Not a SELECT query")
	}

	fromIdx := strings.Index(lowerQuery, "from")
	orderIdx := strings.Index(lowerQuery, "order by")
	groupIdx := strings.Index(lowerQuery, "group by")
	limitIdx := strings.Index(lowerQuery, "limit")

	if fromIdx == -1 {
		return nil, fmt.Errorf("❌ Missing FROM clause")
	}

	// Parse LIMIT
	limit := 0
	limitStart := len(query)
	if limitIdx != -1 {
		limitPart := strings.TrimSpace(query[limitIdx+5:])
		fmt.Sscanf(limitPart, "%d", &limit)
		limitStart = limitIdx
	}

	// Parse GROUP BY correctly
	groupByFields := []string{}
	groupEnd := limitStart
	if orderIdx != -1 && orderIdx < limitStart && orderIdx > groupIdx {
		groupEnd = orderIdx
	}
	if groupIdx != -1 && groupIdx < limitStart {
		groupPart := strings.TrimSpace(query[groupIdx+8 : groupEnd])
		parts := strings.Split(groupPart, ",")
		for _, f := range parts {
			groupByFields = append(groupByFields, strings.ToLower(strings.TrimSpace(f)))
		}
		fmt.Printf("DEBUG Parsed GroupBy fields: %+v\n", groupByFields)
	}

	// Parse ORDER BY (unchanged)
	orderByFields := []ast.OrderByField{}
	orderEnd := len(query)
	if orderIdx != -1 && orderIdx < limitStart {
		if limitIdx != -1 && limitIdx > orderIdx {
			orderEnd = limitIdx
		}
		orderPart := strings.TrimSpace(query[orderIdx+8 : orderEnd])
		clauses := strings.Split(orderPart, ",")
		for _, clause := range clauses {
			parts := strings.Fields(strings.TrimSpace(clause))
			if len(parts) == 0 {
				continue
			}
			field := strings.ToLower(parts[0])
			desc := len(parts) > 1 && strings.ToLower(parts[1]) == "desc"
			orderByFields = append(orderByFields, ast.OrderByField{Field: field, Desc: desc})
		}
	}

	// Determine base query for SELECT ... FROM ...
	baseQueryEnd := len(query)
	if groupIdx != -1 {
		baseQueryEnd = groupIdx
	} else if orderIdx != -1 {
		baseQueryEnd = orderIdx
	} else if limitIdx != -1 {
		baseQueryEnd = limitIdx
	}
	baseQuery := query[:baseQueryEnd]

	selectAST, err := parseCoreSelect(baseQuery)
	if err != nil {
		return nil, err
	}

	selectAST.GroupBy = groupByFields
	selectAST.OrderBy = orderByFields
	selectAST.Limit = limit

	return selectAST, nil
}

func parseCoreSelect(query string) (*ast.SelectQueryAST, error) {
	lowerQuery := strings.ToLower(query)
	fromIdx := strings.Index(lowerQuery, "from")
	if fromIdx == -1 {
		return nil, fmt.Errorf("❌ Missing FROM clause")
	}

	selectPart := strings.TrimSpace(query[6:fromIdx])
	fromRest := strings.TrimSpace(query[fromIdx+4:])

	fields := []string{}
	aggregates := []ast.AggregateFunc{}

	selectParts := strings.Split(selectPart, ",")
	for _, part := range selectParts {
		part = strings.TrimSpace(part)
		if matches := aggregateRegex.FindStringSubmatch(part); len(matches) == 3 {
			funcName := strings.ToUpper(matches[1])
			fieldName := matches[2]
			aggregates = append(aggregates, ast.AggregateFunc{Func: funcName, Field: fieldName})
		} else {
			fields = append(fields, strings.ToLower(part))
		}
	}

	tableName, conditions, err := parseFromWhere(fromRest)
	if err != nil {
		return nil, err
	}

	return &ast.SelectQueryAST{
		Fields:     fields,
		Table:      tableName,
		Conditions: conditions,
		Aggregates: aggregates,
	}, nil
}

/*func parseCoreSelect(query string) (*ast.SelectQueryAST, error) {
	lowerQuery := strings.ToLower(query)
	fromIdx := strings.Index(lowerQuery, "from")
	if fromIdx == -1 {
		return nil, fmt.Errorf("❌ Missing FROM clause")
	}

	selectPart := strings.TrimSpace(query[6:fromIdx])
	fromRest := strings.TrimSpace(query[fromIdx+4:])

	fields := []string{}
	aggregates := []ast.AggregateFunc{}

	selectParts := strings.Split(selectPart, ",")
	for _, part := range selectParts {
		part = strings.TrimSpace(part)

		// Check for aggregate
		if matches := aggregateRegex.FindStringSubmatch(part); len(matches) == 3 {
			funcName := strings.ToUpper(matches[1])
			fieldName := matches[2]

			aggregates = append(aggregates, ast.AggregateFunc{Func: funcName, Field: fieldName})
			fields = append(fields, fmt.Sprintf("%s(%s)", funcName, fieldName)) // ADD raw aggregate string to fields
		} else {
			fields = append(fields, strings.ToLower(part))
		}
	}

	tableName, conditions, err := parseFromWhere(fromRest)
	if err != nil {
		return nil, err
	}

	return &ast.SelectQueryAST{
		Fields:     fields,
		Table:      tableName,
		Conditions: conditions,
		Aggregates: aggregates,
	}, nil
}
*/
