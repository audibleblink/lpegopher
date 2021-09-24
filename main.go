package main

import (
	"github.com/alexflint/go-arg"
)

var cli = arg.MustParse(&args)

func main() {
	switch {
	case args.Collect != nil:
		handleCollect(cli, args)
	case args.Process != nil:
		err := handleProcess(cli, args)
		if err != nil {
			// panic(err)
		}
	}
}
