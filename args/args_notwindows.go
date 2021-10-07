//go:build !windows

package args

type dummy struct{}
type collectCmd struct {
	Tasks    *dummy `arg:"subcommand" help:"Only available on Windows"`
	Services *dummy `arg:"subcommand" help:"Only available on Windows"`
	Exes     *dummy `arg:"subcommand" help:"Only available on Windows"`
	Dlls     *dummy `arg:"subcommand" help:"Only available on Windows"`
}
