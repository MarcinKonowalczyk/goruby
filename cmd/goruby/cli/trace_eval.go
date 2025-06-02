package cli

import (
	"flag"
	"fmt"
)

var trace_eval string

func initTraceEval() {
	help := "trace evaluation of Ruby code. Options: off, on, only, on-no-messages, only-no-messages"
	flag.StringVar(&trace_eval, "trace-eval", "off", help)
}

type TraceEval string

const (
	TraceEval_Off           TraceEval = "off"
	TraceEval_On            TraceEval = "on"
	TraceEval_On_NoMessages TraceEval = "on-no-messages"
)

func parseTraceEval() (TraceEval, error) {
	switch trace_eval {
	case "off", "false":
		return TraceEval_Off, nil
	case "on", "true":
		return TraceEval_On, nil
	case "on-no-messages":
		return TraceEval_On_NoMessages, nil
	default:
		return TraceEval_Off, fmt.Errorf("invalid value for -trace-eval: %s. Valid values are: off, on, only", trace_eval)
	}
}

func (t TraceEval) On() bool {
	switch t {
	case TraceEval_On, TraceEval_On_NoMessages:
		return true
	default:
		return false
	}
}

func (t TraceEval) NoMessages() bool {
	switch t {
	case TraceEval_On_NoMessages:
		return true
	default:
		return false
	}
}
