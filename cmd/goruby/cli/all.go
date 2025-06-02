package cli

func Init() {
	initTraceParse()
	initTraceEval()
	initNoPrint()
}

type Flags struct {
	TraceParse TraceParse
	TraceEval  TraceEval
	NoPrint    bool
}

func Parse() (f Flags, err error) {
	traceParse, err := parseTraceParse()
	if err != nil {
		return f, err
	}

	traceEval, err := parseTraceEval()
	if err != nil {
		return f, err
	}

	noPrint := parseNoPrint()

	return Flags{
		TraceParse: traceParse,
		TraceEval:  traceEval,
		NoPrint:    noPrint,
	}, nil
}
