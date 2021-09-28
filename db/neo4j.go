package db

import (
	"github.com/audibleblink/pegopher/args"
	"github.com/mindstand/gogm/v2"
)

var argv = args.ArgType{}

func Session() (sess gogm.SessionV2, err error) {

	sessConf := gogm.SessionConfig{
		AccessMode:   gogm.AccessModeWrite,
		DatabaseName: argv.Process.Database,
	}

	return gogm.G().NewSessionV2(sessConf)
}
