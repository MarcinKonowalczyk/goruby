package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	transformer_pipeline "github.com/MarcinKonowalczyk/goruby/pipelines/transformer"
)

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
