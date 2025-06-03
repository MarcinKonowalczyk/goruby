package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

type Printer interface {
	io.Writer
	PrintNode(ast.Node)
	Logf(...string)
	String() string
}

func NewPrinter(prefix string) Printer {
	return &printer{
		out:            &strings.Builder{},
		comment_prefix: prefix,
	}
}

func (g *printer) Write(p []byte) (n int, err error) {
	return g.out.Write(p)
}

var _ io.Writer = &printer{}

type printer struct {
	out            *strings.Builder
	comment_prefix string
}

func (g *printer) Println(args ...interface{}) {
	g.out.WriteString(fmt.Sprintln(args...))
}

func (g *printer) Printf(format string, args ...interface{}) {
	g.out.WriteString(fmt.Sprintf(format, args...))
}

func (p *printer) PrintNode(node ast.Node) {
	switch node := (node).(type) {
	case *ast.Program:
		for _, statement := range node.Statements {
			if statement == nil {
				continue
			}
			p.PrintNode(statement.(ast.Node))
		}
	case *ast.Comment:
		p.Println(node.Code())
	case *ast.ExpressionStatement:
		p.PrintNode(node.Expression)
	case *ast.FunctionLiteral:
		p.Logf(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()
	case *ast.Assignment:
		p.Logf(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()

	case *ast.IndexExpression:
		p.Logf(fmt.Sprintf(" %T", node))
		p.Println(node.Code())

	case *ast.ContextCallExpression:
		p.Logf(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()
	case *ast.IntegerLiteral:
		p.Logf(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()
	case *ast.ConditionalExpression:
		p.Logf(fmt.Sprintf("%T", node))
		p.Println(node.Code())
		p.Println()

	default:
		panic(fmt.Sprintf("GRGR print does not yet know how to print %T", node))
	}
}

var _ Printer = &printer{}

func (p *printer) Logf(content ...string) {
	if len(content) == 0 {
		return
	}

	var rest []any
	if len(content) > 1 {
		for i := 1; i < len(content); i++ {
			rest = append(rest, content[i])
		}
	} else {
		rest = []any{}
	}

	msg := fmt.Sprintf(content[0], rest...)
	msg = strings.TrimPrefix(msg, " ")
	msg = strings.TrimSuffix(msg, "\n")
	msg = " " + msg
	comment_line := &ast.Comment{Value: p.comment_prefix + msg}
	p.Println(comment_line.Code())
}

func (p *printer) String() string {
	return p.out.String()
}
