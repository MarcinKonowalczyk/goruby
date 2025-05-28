package printer

import (
	"fmt"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

type Printer interface {
	PrintNode(ast.Node)
	String() string
}

func NewPrinter(prefix string) Printer {
	return &printer{
		out:            &strings.Builder{},
		comment_prefix: prefix,
	}
}

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
		p.PrintFakeComment(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()
	case *ast.Assignment:
		p.PrintFakeComment(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()

	case *ast.IndexExpression:
		p.PrintFakeComment(fmt.Sprintf(" %T", node))
		p.Println(node.Code())

	case *ast.ContextCallExpression:
		p.PrintFakeComment(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()
	case *ast.IntegerLiteral:
		p.PrintFakeComment(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()
	case *ast.ConditionalExpression:
		p.PrintFakeComment(fmt.Sprintf(" %T", node))
		p.Println(node.Code())
		p.Println()

	default:
		panic(fmt.Sprintf("GRGR print does not yet know how to print %T", node))
	}
}

var _ Printer = &printer{}

func (p *printer) PrintFakeComment(content ...string) {
	if len(content) == 0 {
		return
	}
	for _, line := range content {
		comment_line := &ast.Comment{Value: p.comment_prefix + line}
		p.Println(comment_line.Code())
	}
}

func (p *printer) String() string {
	return p.out.String()
}
