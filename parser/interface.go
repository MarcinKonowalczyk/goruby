package parser

import (
	"context"
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

type Parser interface {
	Parse(src string) (*ast.Program, error)
	ParseCtx(ctx context.Context, src string) (*ast.Program, error)
}

func NewParser() Parser {
	var p parser
	p.init()
	return &p
}

func (p *parser) Parse(src string) (*ast.Program, error) {
	return p.parse(context.Background(), src)
}

func (p *parser) ParseCtx(ctx context.Context, src string) (*ast.Program, error) {
	return p.parse(ctx, src)
}

var _ Parser = (*parser)(nil)

// Parse source
func Parse(src string) (*ast.Program, error) {
	p := NewParser()
	return p.Parse(src)
}

func ParseCtx(ctx context.Context, src string) (*ast.Program, error) {
	p := NewParser()
	return p.ParseCtx(ctx, src)
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
