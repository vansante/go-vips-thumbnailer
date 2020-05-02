package thumbnailer

var (
	logger Logger = NoLogger{}
)

func SetLogger(newLogger Logger) {
	logger = newLogger
}

// Logger is an optional interface used for outputting debug logging
type Logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type NoLogger struct{}

func (l NoLogger) Debugf(format string, args ...interface{}) {
}

func (l NoLogger) Errorf(format string, args ...interface{}) {
}
