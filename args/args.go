package args

var Args ArgType

type neoConnection struct {
	Host     string `arg:"-d,env:NEO_HOST" default:"localhost"`
	Password string `arg:"-w,env:NEO_PASSWORD" default:"neo4j"`
	Port     int    `arg:"-p,env:NEO_PORT" default:"7687"`
	Username string `arg:"-u,env:NEO_USER" default:"neo4j"`
	Protocol string `arg:"-t,env:NEO_PROTO" default:"bolt+s"`
	Database string `arg:"-t,env:NEO_DBNAME" default:"neo4j"`
}

type ArgType struct {
	Collect     *collectCmd   `arg:"subcommand" help:"Collect necsesary data"`
	PostProcess *processCmd   `arg:"subcommand" help:"Run Post-Processing tasks and populate neo4j"`
	GetSystem   *getSystemCmd `arg:"subcommand" help:"Utility for acquiring SYSTEM"`

	Debug bool `arg:"-v" help:"verbose output" default:"false"`

	neoConnection
}

type getSystemCmd struct {
	PID int `help:"Process PID that's running as system (defaults to winlogon.exe)"`
}

type processCmd struct {
	Runners string `arg:"-r" help:"Path to collected Runners json"`
	Drop    bool   `help:"drop the database before processing" default:"false"`
}
