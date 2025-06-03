package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/cmd/goruby/cli"
	it "github.com/MarcinKonowalczyk/goruby/pipelines/interpreter"
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
var cpuprofile string = ""

func main() {
	cli.Init()
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.Var(&onelineScripts, "e", "one line of script. Several -e's allowed. Omit [programfile]")
	version := flag.Bool("version", false, "print version")
	flag.Parse()
	flags, err := cli.Parse()
	if err != nil {
		log.Fatal("Error parsing flags: ", err)
	}

	if *version {
		fmt.Println("goruby 0.1.0")
		os.Exit(0)
	}
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal("could not close CPU profile: ", err)
			}
		}()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	args := flag.Args()
	if len(args) == 0 && len(onelineScripts) == 0 {
		log.Println("No program files specified")
		os.Exit(1)
	}

	// set up the interpreter
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	if flags.NoPrint {
		stdout = nil
	}

	interpreter := it.NewInterpreter(stdin, stdout, stderr, args[1:])

	interpreter.SetOptions(func(opts *it.Options) {
		opts.TraceParse = flags.TraceParse.On()
		opts.PrintParseMessages = !flags.TraceParse.NoMessages()

		opts.TraceEval = flags.TraceEval.On()
		opts.PrintEvalMessages = !flags.TraceEval.NoMessages()
	})

	// if we have oneline scripts, interpret them
	if len(onelineScripts) != 0 {
		input := strings.Join(onelineScripts, "\n")
		_, err := interpreter.InterpretCode(input)
		if err != nil {
			fmt.Printf("%T:%v\n", err, err)
			os.Exit(1)
		}
		return
	}

	if len(args) == 0 {
		log.Fatalf("No program files specified")
	}
	fileBytes, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatalf("Error while opening program file: %T:%v\n", err, err)
	}

	if flags.TraceParse.Only() {
		// we can only parse the code, not interpret it
		_, err = interpreter.ParseCode(string(fileBytes))
	} else {
		_, err = interpreter.InterpretCode(string(fileBytes))
	}

	if err != nil {
		fmt.Printf("%T : %v\n", err, err)
		os.Exit(1)
	}
}
