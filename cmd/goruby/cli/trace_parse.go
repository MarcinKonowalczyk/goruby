package cli

import (
	"flag"
	"fmt"
)

var trace_parse string

func initTraceParse() {
	flag.StringVar(&trace_parse, "trace-parse", "off", "trace parsing of Ruby code. Options: off, on, only")
}

type TraceParse string

const (
	TraceParse_Off  TraceParse = "off"
	TraceParse_On   TraceParse = "on"
	TraceParse_Only TraceParse = "only"
)

func parseTraceParse() (TraceParse, error) {
	switch trace_parse {
	case "off", "false":
		return TraceParse_Off, nil
	case "on", "true":
		return TraceParse_On, nil
	case "only":
		return TraceParse_Only, nil
	default:
		return TraceParse_Off, fmt.Errorf("invalid value for -trace-parse: %s. Valid values are: off, on, only", trace_parse)
	}
}

func (t TraceParse) Enabled() bool {
	switch t {
	case TraceParse_On, TraceParse_Only:
		return true
	default:
		return false
	}
}
