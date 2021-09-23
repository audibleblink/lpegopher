//go:build !windows

package main

import (
	"github.com/alexflint/go-arg"
)

func handleCollect(cli *arg.Parser, a argType) {
	cli.Fail("collect functionality is only available on Windows")
}

func handleProcess(cli *arg.Parser, a argType) {

}
