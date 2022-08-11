package mint

import "log"

const (
	logPrefix = "artwork/mint: "
)

func init() {
	var (
		logErr = log.Default
		log    = log.New(os.StdOut, logPrefix, log.LstdFlags)
	)
	logErr.SetPrefix(logPrefix)
}
