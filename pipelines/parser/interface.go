package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

func readSource(filename string, src interface{}) (string, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return s, nil
		case []byte:
			return string(s), nil
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.String(), nil
			}
		case io.Reader:
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, s); err != nil {
				return "", err
			}
			return buf.String(), nil
		}
		return "", fmt.Errorf("unsupported type for src: %T", src)
	}
	// just read the file
	file_src, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read source file %q: %w", filename, err)
	}
	return string(file_src), nil

}

func ParseFile(filename string, trace_parse bool) (*ast.Program, trace.Tracer, error) {
	// get source
	src, err := readSource(filename, nil)
	if err != nil {
		return nil, nil, err
	}

	p := parser.NewParser(src, trace_parse)

	program, err := p.Parse()
	tracer := p.Tracer()
	if tracer != nil {
		tracer.Done()
	}
	return program, tracer, err
}

// // ParseExprFrom is a convenience function for parsing an expression.
// // The arguments have the same meaning as for ParseFile, but the source must
// // be a valid Go (type or value) expression.
// func ParseExprFrom(filename string, src interface{}) (ast.Expression, error) {
// 	// get source
// 	text, err := readSource(filename, src)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var p parser
// 	p.init(filename, text, false)

// 	program, err := p.ParseProgram()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(program.Statements) == 0 {
// 		return nil, fmt.Errorf("source did not contain any expressions to parse")
// 	}

// 	for _, stmt := range program.Statements {
// 		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
// 			return exprStmt.Expression, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("source only contains statements")
// }

// // ParseExpr is a convenience function for obtaining the AST of an expression x.
// // The position information recorded in the AST is undefined. The filename used
// // in error messages is the empty string.
// func ParseExpr(x string) (ast.Expression, error) {
// 	return ParseExprFrom("", []byte(x))
// }
