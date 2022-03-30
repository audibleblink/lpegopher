package args

var Args ArgType

type ArgType struct {
	GetSystem *getSystemCmd `arg:"subcommand" help:"Utility for acquiring SYSTEM before collection"`
	Collect   *collectCmd   `arg:"subcommand" help:"Collect Windows PE and Runner data"`
	Process   *processCmd   `arg:"subcommand" help:"Run Post-Processing tasks and populate neo4j"`

	Debug   bool `arg:"-v" help:"verbose output" default:"false"`
	NoColor bool `arg:"--nocolor" help:"Disable colored output" default:"false"`
}

type getSystemCmd struct {
	PID int `help:"Process PID that's running as system (defaults to winlogon.exe)"`
}

type processCmd struct {
	Drop bool   `help:"drop the database before processing" default:"false"`
	HTTP string `help:"stage files from cwd instead of needing to upload to /import"`

	Username string `arg:"--user,env:NEO_USER" default:"neo4j" placeholder:"<user>"`
	Password string `arg:"--pass,env:NEO_PASSWORD" default:"neo4j" placeholder:"<pass>"`
	Host     string `arg:"env:NEO_HOST" default:"localhost" placeholder:"<host>"`
	Port     int    `arg:"env:NEO_PORT" default:"7687" placeholder:"<port>"`
	Database string `arg:"--db,env:NEO_DBNAME" default:"neo4j" placeholder:"<dbname>"`
	Protocol string `arg:"--proto,env:NEO_PROTO" default:"bolt" placeholder:"<proto>"`
}
