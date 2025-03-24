// ast/insert_ast.go
package ast

type InsertQueryAST struct {
	Table   string
	Columns []string
	Values  []string
}
