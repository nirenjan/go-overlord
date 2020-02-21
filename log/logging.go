// Package log provides primitives to display log messages to stderr.
// This is used internally within Overlord to display messages at various
// levels of intensity
package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
)

// Type Level represents a logging level. Messages below the set logging level
// are suppressed.
type Level int

// The logging module supports a number of different levels
const (
	// Log messages at Debug level. This is used only for debugging
	DEBUG Level = iota

	// Log messages at Info level. This is typically used for notifying
	// the user that something happened.
	INFO

	// Log messages at Warning level. This is used when something has
	// not worked as expected, but there is no action to be taken. This
	// is also the default logging level.
	WARNING

	// Log messages at Error level. This is used when something has gone
	// wrong and the user may have to perform some action to resolve it.
	ERROR

	// Log fatal messages. This is used when things have gone really wrong
	// and the program must terminate.
	FATAL
)

// String converts the level integer into its string representation
func (l Level) String() string {
	return [...]string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}[l]
}

// System log level, default is WARNING
var loglevel = WARNING

func init() {
	l, ok := os.LookupEnv("OVERLORD_LOG")
	if ok {
		switch l {
		case "DEBUG", "4", "5", "6", "7", "8", "9":
			loglevel = DEBUG

		case "INFO", "3":
			loglevel = INFO

		case "WARNING", "2":
			loglevel = WARNING

		case "ERROR", "1":
			loglevel = ERROR

		case "FATAL", "0":
			loglevel = FATAL

		default:
			// Anything else is not recognized, leave it at WARNING
		}
	}
}

// SetLevel sets the logging level. If this is never called, the logging level
// is set to WARNING
func SetLevel(level Level) {
	loglevel = level
}

// GetLevel retrieves the current logging level
func GetLevel() Level {
	return loglevel
}

// System output stream, must implement io.Writer
var output_stream io.Writer = os.Stderr

// SetOutput sets the output location for log writes. It defaults to os.Stderr
func SetOutput(output io.Writer) {
	output_stream = output
}

// format_log formats the log message and prints it to stderr
func format_log(level Level, args ...interface{}) {
	if loglevel <= level {
		pargs := append([]interface{}{level.String()}, args...)
		output_stream.Write([]byte(fmt.Sprintln(pargs...)))
	}
}

// Debug logs a message at DEBUG level.
// This will also capture the calling file and line number
func Debug(args ...interface{}) {
	if loglevel > DEBUG {
		return
	}

	// Get the file name and line number of the caller
	_, file, line, ok := runtime.Caller(1)

	var pargs []interface{}
	if ok {
		file = path.Base(file)
		prefix := fmt.Sprintf("%v:%v -", file, line)
		pargs = append([]interface{}{prefix}, args...)
	} else {
		pargs = args
	}
	format_log(DEBUG, pargs...)
}

// Info logs a message at Info level
func Info(args ...interface{}) {
	format_log(INFO, args...)
}

// Warning logs a message at Warning level
func Warning(args ...interface{}) {
	format_log(WARNING, args...)
}

// Error logs a message at Error level
func Error(args ...interface{}) {
	format_log(ERROR, args...)
}

// Fatal logs a message at Fatal level and exit the program
func Fatal(args ...interface{}) {
	format_log(FATAL, args...)
	os.Exit(1)
}
