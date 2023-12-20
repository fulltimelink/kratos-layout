package boot

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

type bootLogger struct {
	Id string
}

func NewBootLogger(id string) *bootLogger {
	return &bootLogger{Id: id}
}

func (l *bootLogger) Run() log.Logger {
	return log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", l.Id,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
}
