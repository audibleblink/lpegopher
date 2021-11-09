package args

type collectCmd struct {
	Runners *collectRunnersCmd `arg:"subcommand"`
	PEs     *collectPECmd      `arg:"subcommand"`
	All     bool               `arg:"--all" help:"Collect everything into files: ./{runners,pes}.json"`
}

type collectRunnersCmd struct {
	File string `arg:"--outfile,-o" help:"Output file name" default:"stdOut"`
}

type collectPECmd struct {
	File string `arg:"--outfile,-o" help:"Output file name" default:"stdOut"`
	Path string `arg:"required" help:"Directory whence recursive searching begins"`
}
