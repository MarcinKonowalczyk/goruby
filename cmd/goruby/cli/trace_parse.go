package cli

import (
	"flag"
	"fmt"
)

var trace_parse string

func initTraceParse() {
	help := "trace parsing of Ruby code. Options: off, on, only, on-no-messages, only-no-messages"
	flag.StringVar(&trace_parse, "trace-parse", "off", help)
}

type TraceParse string

const (
	TraceParse_Off             TraceParse = "off"
	TraceParse_On              TraceParse = "on"
	TraceParse_On_NoMessages   TraceParse = "on-no-messages"
	TraceParse_Only            TraceParse = "only"
	TraceParse_Only_NoMessages TraceParse = "only-no-messages"
)

func parseTraceParse() (TraceParse, error) {
	switch trace_parse {
	case "off", "false":
		return TraceParse_Off, nil
	case "on", "true":
		return TraceParse_On, nil
	case "only":
		return TraceParse_Only, nil
	case "on-no-messages":
		return TraceParse_On_NoMessages, nil
	case "only-no-messages":
		return TraceParse_Only_NoMessages, nil
	default:
		return TraceParse_Off, fmt.Errorf("invalid value for -trace-parse: %s. Valid values are: off, on, only", trace_parse)
	}
}

func (t TraceParse) On() bool {
	switch t {
	case TraceParse_On, TraceParse_Only, TraceParse_On_NoMessages, TraceParse_Only_NoMessages:
		return true
	default:
		return false
	}
}

func (t TraceParse) NoMessages() bool {
	switch t {
	case TraceParse_On_NoMessages, TraceParse_Only_NoMessages:
		return true
	default:
		return false
	}
}

func (t TraceParse) Only() bool {
	switch t {
	case TraceParse_Only, TraceParse_Only_NoMessages:
		return true
	default:
		return false
	}
}
