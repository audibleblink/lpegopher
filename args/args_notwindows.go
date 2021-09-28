//go:build !windows

package args

type ArgType struct {
	Collect *collectCmd `arg:"subcommand" help:"Only available on Windows"`
	Process *processCmd `arg:"subcommand" help:"Process data and populate neo4j"`
}

// var args = argType{}

type dummy struct{}
type collectCmd struct {
	Tasks    *dummy `arg:"subcommand" help:"Only available on Windows"`
	Services *dummy `arg:"subcommand" help:"Only available on Windows"`
	Exes     *dummy `arg:"subcommand" help:"Only available on Windows"`
	Dlls     *dummy `arg:"subcommand" help:"Only available on Windows"`
}
