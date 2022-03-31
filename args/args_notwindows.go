//go:build !windows

package args

type collectCmd struct {
	Root string `arg:"positional,required" help:"This command is only available on Windows"`
}
