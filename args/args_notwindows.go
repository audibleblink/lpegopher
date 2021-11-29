//go:build !windows

package args

type collectCmd struct {
	Path string `arg:"required" help:"This command is only available on Windows"`
}
