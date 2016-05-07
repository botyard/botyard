package log

import (
	"github.com/go-kit/kit/log"

	stdlog "log"
	"os"
)

var (
	Logger log.Logger
)

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.NewContext(Logger).With("ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	stdlog.SetOutput(log.NewStdlibAdapter(Logger))
}
