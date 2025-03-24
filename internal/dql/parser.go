package dql

import (
	"dagenie/internal/dql/ast"
	"dagenie/internal/dql/parser"
)

// ---------------------- INSERT PARSER ----------------------

type DQLInsert struct {
	DAGID   string
	ID      string
	Name    string
	Payload string
	Status  string
}

// ParseInsert parses a fixed-format INSERT INTO dag(...) VALUES(...) query.
func ParseInsert(query string) (*ast.InsertQueryAST, error) {
	return parser.ParseInsertToAST(query)
}

// ---------------------- SELECT PARSER via parser module ----------------------

// ParseSelectToAST delegates SELECT parsing to internal parser
func ParseSelect(query string) (*ast.SelectQueryAST, error) {
	return parser.ParseSelectToAST(query)
}
