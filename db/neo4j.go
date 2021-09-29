package db

import (
	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/mindstand/gogm/v2"
)

var (
	argv = args.ArgType{}
	_    = arg.MustParse(&argv)
)

var cachedSession gogm.SessionV2

func Session() (gogm.SessionV2, error) {
	var err error
	if cachedSession != nil {
		return cachedSession, err
	}

	config := gogm.SessionConfig{
		AccessMode:   gogm.AccessModeWrite,
		DatabaseName: argv.Process.Database,
	}

	cachedSession, err = gogm.G().NewSessionV2(config)
	return cachedSession, err

}
