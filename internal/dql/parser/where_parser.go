package parser

import (
	"dagenie/internal/dql/ast"
	"fmt"
	"regexp"
	"strings"
)

// ------------------------------------------
// WHERE Parser State
// ------------------------------------------
type whereParser struct {
	tokens []string
	pos    int
}

// Define this at the top of your file (outside any function)
var whereTokenPattern = regexp.MustCompile(`\s*(\(|\)|\bAND\b|\bOR\b|\bNOT\b|<=|>=|!=|=|<|>|\bLIKE\b|'[^']*'|"[^"]*"|[^\s()=<>!]+)\s*`)

// Tokenize WHERE clause into individual tokens
func tokenizeWhere(input string) []string {
	matches := whereTokenPattern.FindAllStringSubmatch(input, -1)
	var tokens []string
	for _, match := range matches {
		token := strings.TrimSpace(match[1])
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func ParseWhereClause(fromRest string) (string, ast.LogicalNode, error) {
	fromRest = strings.TrimSpace(fromRest)
	lower := strings.ToLower(fromRest)
	whereIdx := strings.Index(lower, "where")

	var tableName, wherePart string
	if whereIdx != -1 {
		tableName = strings.TrimSpace(fromRest[:whereIdx])
		wherePart = strings.TrimSpace(fromRest[whereIdx+5:])
	} else {
		tableName = fromRest
	}

	if tableName == "" {
		return "", nil, fmt.Errorf("❌ Missing table name after FROM")
	}

	if wherePart == "" {
		return strings.ToLower(tableName), nil, nil
	}

	tokens := tokenizeWhere(wherePart)
	fmt.Printf("[DEBUG] Tokens: %v\n", tokens)

	logicalExpr, err := parseWhereTokens(tokens)
	if err != nil {
		return "", nil, err
	}
	fmt.Printf("[DEBUG] Parsed LogicalExpr: %+v\n", logicalExpr)

	return strings.ToLower(tableName), logicalExpr, nil
}

func parseWhereTokens(tokens []string) (ast.LogicalNode, error) {
	p := &whereParser{tokens: tokens}
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.pos < len(p.tokens) {
		return nil, fmt.Errorf("❌ Unexpected token after WHERE clause: %s", p.peek())
	}

	return expr, nil
}

func (p *whereParser) peek() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return strings.ToUpper(p.tokens[p.pos])
}

func (p *whereParser) consume() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	token := p.tokens[p.pos]
	p.pos++
	return token
}

// parseExpression: handles OR level
func (p *whereParser) parseExpression() (ast.LogicalNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.peek() == "OR" {
		p.consume()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.OrNode{Left: left, Right: right}
	}
	return left, nil
}

// parseAnd: handles AND level
func (p *whereParser) parseAnd() (ast.LogicalNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}

	for p.peek() == "AND" {
		p.consume()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &ast.AndNode{Left: left, Right: right}
	}
	return left, nil
}

// parseNot: handles NOT or passes to parseAtom
func (p *whereParser) parseNot() (ast.LogicalNode, error) {
	if p.peek() == "NOT" {
		p.consume()
		expr, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		return &ast.NotNode{Expr: expr}, nil
	}
	return p.parseAtom()
}

// parseAtom: parses (expr) or condition
func (p *whereParser) parseAtom() (ast.LogicalNode, error) {
	token := p.peek()

	if token == "(" {
		p.consume()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.consume() != ")" {
			return nil, fmt.Errorf("❌ Expected ')' after expression")
		}
		return expr, nil
	}

	// Parse condition: field = value
	field := p.consume()
	if field == "" {
		return nil, fmt.Errorf("❌ Expected field in condition")
	}

	operator := p.consume()
	if operator != "=" {
		return nil, fmt.Errorf("❌ Only '=' operator supported (got '%s')", operator)
	}

	val := p.consume()
	if val == "" {
		return nil, fmt.Errorf("❌ Missing value after '=' for field '%s'", field)
	}
	val = strings.Trim(val, `"'`) // remove quotes

	return &ast.ConditionNode{
		Field:    strings.ToLower(field),
		Operator: "=",
		Value:    val,
	}, nil
}
