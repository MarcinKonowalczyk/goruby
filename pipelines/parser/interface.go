package parser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
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

func ParseFile(ctx context.Context, filename string) (*ast.Program, error) {
	// get source
	src, err := readSource(filename, nil)
	if err != nil {
		return nil, err
	}

	p := parser.NewParser()

	program, err := p.ParseCtx(ctx, src)
	return program, err
}
