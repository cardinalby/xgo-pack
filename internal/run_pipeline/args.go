package run_pipeline

import (
	"flag"
	"os"

	flago "github.com/cardinalby/go-struct-flags"
)

type argsStruct struct {
	Root  string `flag:"root" flagUsage:"Path to the project root directory"`
	Cfg   string `flag:"cfg" flagUsage:"Path to the config file relative to the root"`
	Quiet bool   `flag:"quiet" flagUsage:"Do not print any output"`
}

func getArgs(args []string, defaultCfgRelPath string) argsStruct {
	fls := flago.NewFlagSet("xgo-pack", flag.PanicOnError)
	a := argsStruct{
		Cfg:  defaultCfgRelPath,
		Root: getWd(),
	}
	if err := fls.StructVar(&a); err != nil {
		panic(err)
	}
	if err := fls.Parse(args); err != nil {
		panic(err)
	}
	return a
}

func getWd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
