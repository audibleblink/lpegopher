package args

type collectCmd struct {
	Runners *collectRunnersCmd `arg:"subcommand" help:"collect tasks, services"`
	PEs     *collectINodesCmd  `arg:"subcommand" help:"collect PEs, Dirs, and all their Deps and DACLs"`
	All     bool               `help:"Collect everything into files: ./{runners,inodes}.csv"`
	JSON    bool               `help:"Save results as JSON for use in other tools"`
}

type collectRunnersCmd struct {
	File string `arg:"--outfile,-o" help:"Output file name, use - for stdout" default:"runners.csv"`
}

type collectINodesCmd struct {
	File string `arg:"--outfile,-o" help:"Output file name, use - for stdout" default:"inodes.csv"`
	Path string `arg:"required" help:"Directory whence recursive searching begins"`
}
