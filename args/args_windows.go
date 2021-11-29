package args

type collectCmd struct {
	Path string `arg:"required" help:"Directory whence recursive searching begins"`
}
