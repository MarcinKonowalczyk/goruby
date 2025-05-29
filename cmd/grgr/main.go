package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	transformer_pipeline "github.com/MarcinKonowalczyk/goruby/pipelines/transformer"
)

// var (
// 	trace_transform bool = false
// )

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Println("No program files specified")
		os.Exit(1)
	}

	src, err := os.ReadFile(args[0])
	if err != nil {
		log.Printf("Error reading program file %s: %v\n", args[0], err)
		os.Exit(1)
	}
	transformed_string, err := transformer_pipeline.Transform(string(src))
	if err != nil {
		log.Printf("Error while transforming program file with pipeline: %T:%v\n", err, err)
		os.Exit(1)
	}

	fmt.Printf("%s", transformed_string)
}

// type grgrOutput struct {
// 	out strings.Builder
// }

// func (g *grgrOutput) Println(args ...interface{}) {
// 	g.out.WriteString(fmt.Sprintln(args...))
// }

// func (g *grgrOutput) Printf(format string, args ...interface{}) {
// 	g.out.WriteString(fmt.Sprintf(format, args...))
// }

// func (g *grgrOutput) PrintNode(node ast.Node) {
// 	switch node := (node).(type) {
// 	case *ast.Program:
// 		for _, statement := range node.Statements {
// 			if statement == nil {
// 				continue
// 			}
// 			g.PrintNode(statement.(ast.Node))
// 		}
// 	case *ast.Comment:
// 		g.Println(node.Code())
// 	case *ast.ExpressionStatement:
// 		g.PrintNode(node.Expression)
// 	case *ast.FunctionLiteral:
// 		g.PrintFakeComment(fmt.Sprintf(" %T", node))
// 		g.Println(node.Code())
// 		g.Println()
// 	case *ast.Assignment:
// 		g.PrintFakeComment(fmt.Sprintf(" %T", node))
// 		g.Println(node.Code())
// 		g.Println()

// 	case *ast.IndexExpression:
// 		g.PrintFakeComment(fmt.Sprintf(" %T", node))
// 		g.Println(node.Code())

// 	case *ast.ContextCallExpression:
// 		g.PrintFakeComment(fmt.Sprintf(" %T", node))
// 		g.Println(node.Code())
// 		g.Println()
// 	case *ast.IntegerLiteral:
// 		g.PrintFakeComment(fmt.Sprintf(" %T", node))
// 		g.Println(node.Code())
// 		g.Println()
// 	case *ast.ConditionalExpression:
// 		g.PrintFakeComment(fmt.Sprintf(" %T", node))
// 		g.Println(node.Code())
// 		g.Println()

// 	default:
// 		panic(fmt.Sprintf("GRGR print does not yet know how to print %T", node))
// 	}
// }

// func (g *grgrOutput) PrintFakeComment(content ...string) {
// 	if len(content) == 0 {
// 		return
// 	}
// 	for _, line := range content {
// 		comment_line := &ast.Comment{Value: "GRGR" + line}
// 		g.Println(comment_line.Code())
// 	}
// }
