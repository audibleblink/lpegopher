module github.com/audibleblink/pegopher

go 1.17

require (
	github.com/Microsoft/go-winio v0.5.1
	github.com/alexflint/go-arg v1.4.2
	github.com/audibleblink/concurrent-writer v0.1.0
	github.com/audibleblink/getsystem v0.1.1
	github.com/audibleblink/rpcls v0.0.0-20210822225556-d855a04ad117
	github.com/capnspacehook/taskmaster v0.0.0-20210519235353-1629df7c85e9
	github.com/kgoins/go-winacl/pkg v0.0.0-00010101000000-000000000000
	github.com/minio/highwayhash v1.0.2
	github.com/neo4j/neo4j-go-driver/v4 v4.4.0
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881
	www.velocidex.com/golang/binparsergen v0.1.0
	www.velocidex.com/golang/go-pe v0.1.1-0.20210915141920-02eb5d611e80
)

require (
	github.com/Velocidex/ordereddict v0.0.0-20200723153557-9460a6764ab8 // indirect
	github.com/Velocidex/pkcs7 v0.0.0-20210524015001-8d1eee94a157 // indirect
	github.com/alexflint/go-scalar v1.0.0 // indirect
	github.com/audibleblink/bamflags v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/rickb777/date v1.14.2 // indirect
	github.com/rickb777/plural v1.2.2 // indirect
	golang.org/x/text v0.3.6 // indirect
)

require (
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
)

replace (
	// github.com/kgoins/go-winacl => github.com/audibleblink/go-winacl v0.0.2
	// github.com/kgoins/go-winacl/pkg => github.com/audibleblink/go-winacl v0.0.2
	github.com/kgoins/go-winacl => c:\users\user\code\go-winacl
	github.com/kgoins/go-winacl/pkg => c:\users\user\code\go-winacl
)
