package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

type Log struct {
	Out       chan string // Logger output channel for testing.
	Errors    int         //Non fatal error count.
	Warnings  int         //Warnings count.
	Verbosity int         // 0 => no verbosity; 1 => verbose; 2 => more verbose ...
}

var highlightColor = []color.Attribute{color.FgGreen, color.Bold}
var errorColor = []color.Attribute{color.FgRed, color.Bold}
var warningColor = []color.Attribute{color.FgRed}

// colorize executes a function with color attributes.
func colorize(attributes []color.Attribute, fn func()) {
	defer color.Unset()
	color.Set(attributes...)
	fn()
}

// output prints a line to `out` if `log.verbosity` is greater than equal or
// equal to `verbosity`. If `log.out` is not nil then the line is written to it
// instead of `stdout` (this feature is used for testing purposes).
func (log *Log) output(out io.Writer, verbosity int, format string, v ...any) {
	if log.Verbosity >= verbosity {
		msg := fmt.Sprintf(format, v...)
		if log.Out == nil {
			fmt.Fprintln(out, msg)
		} else {
			log.Out <- msg
		}
	}
}

// logConsole prints a line to stdout.
func (log *Log) Console(format string, v ...any) {
	log.output(os.Stdout, 0, format, v...)
}

// logVerbose prints a line to stdout if `-v` logVerbose option was specified.
func (log *Log) Verbose(format string, v ...any) {
	log.output(os.Stdout, 1, format, v...)
}

// logVerbose2 prints a a line to stdout the `-v` verbose option was specified more
// than once.
func (log *Log) Verbose2(format string, v ...any) {
	colorize(highlightColor, func() {
		log.output(os.Stdout, 2, format, v...)
	})
}

// logColorize prints a colorized line to stdout.
func (log *Log) Colorize(attributes []color.Attribute, format string, v ...any) {
	colorize(attributes, func() {
		log.Console(format, v...)
	})
}

// logHighlight prints a highlighted line to stdout.
func (log *Log) Highlight(format string, v ...any) {
	log.Colorize(highlightColor, format, v...)
}

// logError prints a line to stderr and increments the error count.
func (log *Log) Error(format string, v ...any) {
	colorize(errorColor, func() {
		log.output(os.Stderr, 0, "error: "+format, v...)
	})
	log.Errors++
}

// logWarning prints a line to stdout and increments the warnings count.
func (log *Log) Warning(format string, v ...any) {
	colorize(warningColor, func() {
		log.output(os.Stdout, 0, "warning: "+format, v...)
	})
	log.Warnings++
}
