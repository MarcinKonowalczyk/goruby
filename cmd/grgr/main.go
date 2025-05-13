package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"

	"github.com/MarcinKonowalczyk/goruby/interpreter"
	"github.com/MarcinKonowalczyk/goruby/parser"
)

func main() {
	flag.Parse()
	args := flag.Args()
	interpreter := interpreter.NewInterpreterEx(args[1:])
	_ = interpreter
	if len(args) == 0 {
		log.Println("No program files specified")
		os.Exit(1)
	}
	fileBytes, err := os.ReadFile(args[0])
	if err != nil {
		log.Printf("Error while opening program file: %T:%v\n", err, err)
		os.Exit(1)
	}
	program, err := parser.ParseFile(token.NewFileSet(), args[0], fileBytes)

	if err != nil {
		log.Printf("Error while interpreting program file: %T:%v\n", err, err)
		os.Exit(1)
	}

	for _, statement := range program.Statements {
		fmt.Println(statement)
	}
}
