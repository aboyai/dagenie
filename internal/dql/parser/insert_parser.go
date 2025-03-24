// parser/insert_parser.go
package parser

import (
	"dagenie/internal/dql/ast"
	"fmt"
	"regexp"
	"strings"
)

// ParseInsertToAST parses an INSERT query into InsertQueryAST.
func ParseInsertToAST(query string) (*ast.InsertQueryAST, error) {
	query = strings.TrimSpace(query)
	lowerQuery := strings.ToLower(query)

	if !strings.HasPrefix(lowerQuery, "insert into") {
		return nil, fmt.Errorf("❌ Not an INSERT query")
	}

	// Regex for INSERT INTO dag (col1, col2, ...) VALUES ('val1', 'val2', ...)
	pattern := `(?i)^insert\s+into\s+(\w+)\s*\(([^)]+)\)\s+values\s*\(([^)]+)\)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(query)
	if len(matches) != 4 {
		return nil, fmt.Errorf("❌ Invalid INSERT syntax. Expected: INSERT INTO dag (col1, ...) VALUES ('val1', ...)")
	}

	table := strings.ToLower(matches[1])
	columns := splitCSV(matches[2])
	values := splitCSV(matches[3])

	if len(columns) != len(values) {
		return nil, fmt.Errorf("❌ Column count does not match value count")
	}

	return &ast.InsertQueryAST{
		Table:   table,
		Columns: columns,
		Values:  values,
	}, nil
}

// Helper: Split CSV, trim spaces, strip quotes from values
func splitCSV(input string) []string {
	parts := strings.Split(input, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		trimmed = strings.Trim(trimmed, `"'`)
		result = append(result, trimmed)
	}
	return result
}
