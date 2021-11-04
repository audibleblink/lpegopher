package args

type neoConnection struct {
	Host     string `arg:"-d,env:NEO_HOST" default:"localhost"`
	Password string `arg:"-w,env:NEO_PASSWORD" default:"neo4j"`
	Port     int    `arg:"-p,env:NEO_PORT" default:"7687"`
	Username string `arg:"-u,env:NEO_USER" default:"neo4j"`
	Protocol string `arg:"-t,env:NEO_PROTO" default:"bolt+s"`
	Database string `arg:"-t,env:NEO_DBNAME" default:"neo4j"`
}

type ArgType struct {
	Collect   *collectCmd   `arg:"subcommand" help:"Collect necsesary data"`
	Process   *processCmd   `arg:"subcommand" help:"Process data and populate neo4j"`
	GetSystem *getSystemCmd `arg:"subcommand" help:"Utility for acquiring SYSTEM"`
}

type getSystemCmd struct {
	PID    int  `arg:"required" help:"Process running as system (ex:winlogon.exe)"`
	Self   bool `arg:"" help:"Impersonate SYSTEM in current shell"`
	RunCmd bool `arg:"" help:"Run cmd.exe with duplicated SYSTEM token"`
}

type processCmd struct {
	PEs     string `arg:"-p" help:"Path to collected PEs json" default:"pes.json"`
	Runners string `arg:"-r" help:"Path to collected Runners json" default:"runners.json"`

	neoConnection
}
