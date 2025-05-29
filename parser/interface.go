package parser

import (
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

type Parser interface {
	Parse() (*ast.Program, error)
	Tracer() trace.Tracer
}

func NewParser(src string, trace_parse bool) Parser {
	var p parser
	p.init(src, trace_parse)
	return &p
}

func (p *parser) Tracer() trace.Tracer {
	return p.tracer
}

var _ Parser = (*parser)(nil)

// Parse source
func Parse(src string) (*ast.Program, error) {
	p := NewParser(src, false)
	return p.Parse()
}

// Parse source and return a single statement
// If the source does not contain exactly one statement, an error is returned.
func ParseStatement(src string) (ast.Statement, error) {
	program, err := Parse(src)
	if err != nil {
		return nil, err
	}
	if len(program.Statements) != 1 {
		return nil, fmt.Errorf("source did not contain exactly one statement to parse")
	}
	return program.Statements[0], nil
}

// Parse source and return a single expression
// If the source does not contain exactly one expression, an error is returned.
func ParseExpr(src string) (ast.Expression, error) {
	stmt, err := ParseStatement(src)
	if err != nil {
		return nil, err
	}
	if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
		return exprStmt.Expression, nil
	}
	return nil, fmt.Errorf("source did not contain an expression to parse, got %T", stmt)
}

func MustParse(src string) *ast.Program {
	program, err := Parse(src)
	if err != nil {
		panic(fmt.Sprintf("could not parse source: %s", err))
	}
	return program
}

func MustParseStatement(src string) ast.Statement {
	stmt, err := ParseStatement(src)
	if err != nil {
		panic(fmt.Sprintf("could not parse source: %s", err))
	}
	return stmt
}

func MustParseExpr(src string) ast.Expression {
	expr, err := ParseExpr(src)
	if err != nil {
		panic(fmt.Sprintf("could not parse source: %s", err))
	}
	return expr
}
