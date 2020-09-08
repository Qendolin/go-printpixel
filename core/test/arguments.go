package test

import (
	"flag"
)

func ParseArgs() {
	flag.BoolVar(&Args.Headless, "headless", false, "Configure for headless testing")
	flag.Parse()
}

type Arguments struct {
	Headless bool
}

var Args Arguments
