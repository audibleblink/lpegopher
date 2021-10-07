package logerr

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

var label map[LogLevel]string = map[LogLevel]string{
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARN",
	LogLevelError: "ERROR",
	LogLevelFatal: "FATAL",
}

// OverrideLogger is a dummy logger to kill the gogm logger
type OverrideLogger struct {

	// Level dictates the LogLevel
	Level LogLevel

	// Output destination
	Output io.Writer

	// Exclusive dictates whether _only_ the configured loglevel messages are shown
	// Defaults to printing everything below the configured LogLevel
	// Fatal, Error, Warn, Info, Debug
	Exclusive bool

	// LogWrappedErrors, when enabled, will print the error text,
	// according to level and context, before returning the error
	LogWrappedErrors bool

	// Additional prefix text to add context to log messages
	context  string
	template string
}

var G = DefaultLogger()

func Debug(s string)                       { G.Debug(s) }
func Debugf(s string, vals ...interface{}) { G.Debugf(s, vals) }

func Info(s string)                       { G.Info(s) }
func Infof(s string, vals ...interface{}) { G.Infof(s, vals) }

func Warn(s string)                       { G.Warn(s) }
func Warnf(s string, vals ...interface{}) { G.Warnf(s, vals) }

func Error(s string)                       { G.Error(s) }
func Errorf(s string, vals ...interface{}) { G.Errorf(s, vals) }

func Fatal(s string)                       { G.Fatal(s) }
func Fatalf(s string, vals ...interface{}) { G.Fatalf(s, vals) }

func Wrap(err error) error { return G.Wrap(err) }

func Context(context string) { G.Context(context) }
func ClearContext()          { G.ClearContext() }

func DefaultLogger() *OverrideLogger {
	logger := &OverrideLogger{
		Level:  LogLevelError,
		Output: os.Stderr,
	}
	return logger.Context("")
}

func (d OverrideLogger) SetAsGlobal() {
	G = &d
}

func (d OverrideLogger) ClearContext() {
	d.context = ""
}

func (d OverrideLogger) Context(s string) *OverrideLogger {
	d.context = s
	d.template = "%-8s%s"
	if d.context != "" {
		d.template = fmt.Sprintf("%%-8s%s: %%s", d.context)
	}
	return &d
}

func (d OverrideLogger) Debug(s string) {
	loggerGen(LogLevelDebug, &d)(s)
}

func (d OverrideLogger) Debugf(s string, vals ...interface{}) {
	loggerGenF(LogLevelDebug, &d)(s, vals)
}

func (d OverrideLogger) Info(s string) {
	loggerGen(LogLevelInfo, &d)(s)
}

func (d OverrideLogger) Infof(s string, vals ...interface{}) {
	loggerGenF(LogLevelInfo, &d)(s, vals)
}

func (d OverrideLogger) Warn(s string) {
	loggerGen(LogLevelWarn, &d)(s)
}

func (d OverrideLogger) Warnf(s string, vals ...interface{}) {
	loggerGenF(LogLevelWarn, &d)(s, vals)
}

func (d OverrideLogger) Error(s string) {
	loggerGen(LogLevelError, &d)(s)
}

func (d OverrideLogger) Errorf(s string, vals ...interface{}) {
	loggerGenF(LogLevelError, &d)(s, vals)
}

func (d OverrideLogger) Fatal(s string) {
	loggerGen(LogLevelFatal, &d)(s)
	os.Exit(1)
}

func (d OverrideLogger) Fatalf(s string, vals ...interface{}) {
	loggerGenF(LogLevelFatal, &d)(s, vals)
	os.Exit(1)
}

func (d OverrideLogger) Wrap(err error) error {
	if d.LogWrappedErrors {
		d.Error(err.Error())
	}
	return fmt.Errorf("%s: %w", d.context, err)
}

// generator that return a configured logger
func loggerGen(level LogLevel, l *OverrideLogger) func(string) {
	label := fmt.Sprintf("[%s]", label[level])

	return func(s string) {
		if (l.Level == level && l.Exclusive) || (l.Level <= level && !l.Exclusive) {
			log.Printf(l.template, label, s)
		}
	}
}

// generator that return a configured formatting logger
func loggerGenF(level LogLevel, l *OverrideLogger) func(string, ...interface{}) {
	label := fmt.Sprintf("[%s]", label[level])
	return func(s string, vals ...interface{}) {
		if (l.Level == level && l.Exclusive) || (l.Level <= level && !l.Exclusive) {
			log.Printf(l.template+" %v", label, s, vals)
		}
	}
}

// // generator that return a configured err logger
// func logWithErrGen(level LogLevel, l *OverrideLogger) func(error) error {
// 	label := fmt.Sprintf("[%s]", label[level])
// 	return func(err error) error {
// 		if (l.Level == level && l.Exclusive) || (l.Level <= level && !l.Exclusive) {
// 			log.Printf(l.template, label, err.Error())
// 		}
// 		return l.Wrap(err)
// 	}
// }
