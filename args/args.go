package args

type neoConnection struct {
	Host     string `arg:"-d,env:NEO_HOST" default:"localhost"`
	Password string `arg:"-w,env:NEO_PASSWORD" default:"neo4j"`
	Port     int    `arg:"-p,env:NEO_PORT" default:"7687"`
	Username string `arg:"-u,env:NEO_USER" default:"neo4j"`
	Protocol string `arg:"-t,env:NEO_PROTO" default:"bolt+s"`
	Database string `arg:"-t,env:NEO_DBNAME" default:"neo4j"`
}

type processCmd struct {
	Tasks    *processTasksCmd    `arg:"subcommand"`
	Services *processServicesCmd `arg:"subcommand"`
	Exes     *processExesCmd     `arg:"subcommand"`
	Dlls     *processDllsCmd     `arg:"subcommand"`

	neoConnection
}

type fileIn struct {
	File string `arg:"required,--file,-f" help:"File to process"`
}

type processTasksCmd struct{ fileIn }
type processServicesCmd struct{ fileIn }
type processExesCmd struct{ fileIn }
type processDllsCmd struct{ fileIn }
