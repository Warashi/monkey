package testdata

import _ "embed"

var (
	//go:embed symbols.monkey
	Symbols string
	//go:embed first.monkey
	First string
	//go:embed second.monkey
	Second string
)
