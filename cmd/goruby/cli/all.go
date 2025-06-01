package cli

func Init() {
	initTraceParse()
}

type Flags struct {
	TraceParse TraceParse
}

func Parse() (Flags, error) {
	traceParse, err := parseTraceParse()
	if err != nil {
		return Flags{}, err
	}
	return Flags{
		TraceParse: traceParse,
	}, nil
}
