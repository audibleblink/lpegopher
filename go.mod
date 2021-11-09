module github.com/audibleblink/pegopher

go 1.17

require (
	github.com/alexflint/go-arg v1.4.2
	github.com/capnspacehook/taskmaster v0.0.0-20210519235353-1629df7c85e9
	github.com/mindstand/gogm/v2 v2.2.0
	golang.org/x/sys v0.0.0-20210921065528-437939a70204
)

require (
	github.com/Microsoft/go-winio v0.5.0
	github.com/adam-hanna/arrayOperations v0.2.6 // indirect
	github.com/alexflint/go-scalar v1.0.0 // indirect
	github.com/cornelk/hashmap v1.0.1 // indirect
	github.com/dchest/siphash v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kgoins/go-winacl v0.2.0
	github.com/mindstand/go-cypherdsl v0.2.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/rickb777/date v1.14.2 // indirect
	github.com/rickb777/plural v1.3.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	www.velocidex.com/golang/binparsergen v0.1.0
	www.velocidex.com/golang/go-pe v0.1.1-0.20210915141920-02eb5d611e80
)

require (
	github.com/Velocidex/ordereddict v0.0.0-20200723153557-9460a6764ab8 // indirect
	github.com/Velocidex/pkcs7 v0.0.0-20210524015001-8d1eee94a157 // indirect
	github.com/audibleblink/bamflags v0.2.0 // indirect
	github.com/audibleblink/getsystem v0.1.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/neo4j/neo4j-go-driver/v4 v4.3.3
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/text v0.3.6 // indirect
)

replace (
	github.com/kgoins/go-winacl => github.com/audibleblink/go-winacl v0.0.2
	github.com/kgoins/go-winacl/pkg => github.com/audibleblink/go-winacl v0.0.2
)
