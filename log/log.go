package log

import protocol "github.com/longbridge/openapi-protocol/go"

type Logger = protocol.Logger

var defaultLogger Logger

func init() {
	defaultLogger = &protocol.DefaultLogger{}
}

// SetLogger use to set defaultLogger
func SetLogger(l Logger) {
	if l != nil {
		defaultLogger = l
	}
}

// SetLevel use to modify log level of default logger
func SetLevel(lvl string) {
	defaultLogger.SetLevel(lvl)
}

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func Debugf(f string, args ...interface{}) {
	defaultLogger.Debugf(f, args...)
}

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Infof(msg string, args ...interface{}) {
	defaultLogger.Infof(msg, args...)
}

func Warn(msg string) {
	defaultLogger.Warn(msg)
}

func Warnf(f string, args ...interface{}) {
	defaultLogger.Warnf(f, args...)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Errorf(f string, args ...interface{}) {
	defaultLogger.Errorf(f, args...)
}
