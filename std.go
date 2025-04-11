package printer

import (
	"fmt"
	"os"
)

// globalPrinter is the default Printer instance used by global logging functions.
var globalPrinter = NewPrint(LevelDebug, FlagWithDate|FlagWithGoroutineID|FlagWithColor, os.Stdout, os.Stderr)

// Printf formats and writes a message to the standard output stream using the global printer.
//
// Parameters:
//   - p: string - The format string.
//   - a: ...interface{} - The arguments to format.
func Printf(p string, a ...interface{}) {
	globalPrinter.WriteToStd([]byte(fmt.Sprintf(p, a...)))
}

// Print writes a plain string to the standard output stream using the global printer.
//
// Parameters:
//   - s: string - The message to write.
func Print(s string) {
	globalPrinter.WriteToStd([]byte(s))
}

// PrintError writes an error message to the error output stream using the global printer.
//
// Parameters:
//   - err: error - The error to write. If nil, "<nil>" is printed.
func PrintError(err error) {
	if err == nil {
		globalPrinter.WriteToErr([]byte("<nil>"))
	} else {
		globalPrinter.WriteToErr([]byte(err.Error()))
	}
}

// PrintErrorS writes a plain error string to the error output stream using the global printer.
//
// Parameters:
//   - err: string - The error message to write.
func PrintErrorS(err string) {
	globalPrinter.WriteToErr([]byte(err))
}

// PrintErrorSf formats and writes an error message to the error output stream using the global printer.
//
// Parameters:
//   - err: string - The format string.
//   - args: ...interface{} - The arguments to format.
func PrintErrorSf(err string, args ...interface{}) {
	globalPrinter.WriteToErr([]byte(fmt.Sprintf(err, args...)))
}

// Errorf logs a formatted error message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func Errorf(format string, a ...interface{}) {
	globalPrinter.Errorf(format, a...)
}

// Warnf logs a formatted warning message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func Warnf(format string, a ...interface{}) {
	globalPrinter.Warnf(format, a...)
}

// Infof logs a formatted informational message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func Infof(format string, a ...interface{}) {
	globalPrinter.Infof(format, a...)
}

// Debugf logs a formatted debug message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func Debugf(format string, a ...interface{}) {
	globalPrinter.Debugf(format, a...)
}

// SetLogLevel updates the global printer's log level.
//
// Parameters:
//   - level: int - The new log level to set.
func SetLogLevel(level Levels) {
	globalPrinter.SetLogLevel(level)
}

// GetLogLevel retrieves the current log level of the global printer.
//
// Returns:
//   - int: The current global log level.
func GetLogLevel() Levels {
	return globalPrinter.GetLogLevel()
}
