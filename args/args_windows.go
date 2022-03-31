package args

type collectCmd struct {
	Root string `arg:"positional,required" help:"Directory whence recursive searching begins"`
}
