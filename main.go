package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/interpreter"
	"github.com/pkg/errors"
)

type multiString []string

func (m multiString) String() string {
	out := ""
	for _, s := range m {
		out += s
	}
	return out
}
func (m *multiString) Set(s string) error {
	*m = append(*m, s)
	return nil
}

var onelineScripts multiString
var trace bool

func printError(err error) {
	// fmt.Printf("%v\n", errors.Cause(err))
	fmt.Printf("%T : %v\n", errors.Cause(err), err)
}

func main() {
	flag.Var(&onelineScripts, "e", "one line of script. Several -e's allowed. Omit [programfile]")
	flag.BoolVar(&trace, "trace", false, "trace execution")
	// flag.BoolVar(&debug, "debug", false, "debug execution")
	flag.Parse()
	// fmt.Println("goruby version 0.1.0")
	interpreter := interpreter.NewInterpreter()
	interpreter.Trace = trace
	if len(onelineScripts) != 0 {
		input := strings.Join(onelineScripts, "\n")
		_, err := interpreter.Interpret("", input)
		if err != nil {
			fmt.Printf("%v\n", errors.Cause(err))
			os.Exit(1)
		}
		return
	}
	args := flag.Args()
	if len(args) == 0 {
		log.Println("No program files specified")
		os.Exit(1)
	}
	fileBytes, err := os.ReadFile(args[0])
	if err != nil {
		log.Printf("Error while opening program file: %T:%v\n", err, err)
		os.Exit(1)
	}
	_, err = interpreter.Interpret(args[0], fileBytes)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}
