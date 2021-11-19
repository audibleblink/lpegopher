//go:build !windows

package args

type dummy2 struct{}
type collectCmd struct {
	Tasks    *dummy2 `arg:"subcommand" help:"Only available on Windows"`
	Services *dummy2 `arg:"subcommand" help:"Only available on Windows"`
	PEs      *dummy2 `arg:"subcommand" help:"Only available on Windows"`
}
