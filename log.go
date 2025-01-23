package artwork

import "log"

const (
	logPrefix = "artwork: "
)

var (
	logErr *log.Logger
)

func init() {
	var (
		logErr = log.Default
		logger = log.New(os.StdOut, logPrefix, log.LstdFlags)
	)
	_ = logger
	logErr.SetPrefix(logPrefix)
}
