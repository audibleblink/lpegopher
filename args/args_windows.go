package args


type collectCmd struct {
	Tasks    *collectTaskCmd     `arg:"subcommand"`
	Services *collectServicesCmd `arg:"subcommand"`
	Exes     *collectPECmd       `arg:"subcommand"`
	Dlls     *collectPECmd       `arg:"subcommand"`
	All      bool                `arg:"--all" help:"Collect everything into files: {tasks,services,exes,dlls}.json"`
}

type collectTaskCmd struct{ fileOut }
type collectServicesCmd struct{ fileOut }
type collectPECmd struct {
	fileOut
	Path string `arg:"required" help:"Directory from where recursive searching will begin"`
}

type fileOut struct {
	File string `arg:"--outfile,-o" help:"Output file name" default:"stdOut"`
}
