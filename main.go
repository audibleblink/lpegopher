package main

import (
	"github.com/alexflint/go-arg"
)

var cli = arg.MustParse(&args)

func main() {

	switch cli.Subcommand().(type) {
	case *collectCmd:
		handleCollect(cli, args)
	case processCmd:
		handleProcess(cli, args)
	}
}
