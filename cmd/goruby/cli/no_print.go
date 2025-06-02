package cli

import (
	"flag"
)

var no_print bool

func initNoPrint() {
	help := "supress print and puts from the interpreter. Useful for testing."
	flag.BoolVar(&no_print, "no-output", false, help)
}

func parseNoPrint() bool {
	return no_print
}
