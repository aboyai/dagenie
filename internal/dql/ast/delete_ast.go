package ast

// DeleteQueryAST represents a DELETE ... WHERE ... query
type DeleteQueryAST struct {
	Table      string            // e.g., "dag"
	Conditions map[string]string // WHERE conditions
}
