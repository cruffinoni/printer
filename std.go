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
//   - a: ...any - The arguments to format.
func Printf(p string, a ...any) {
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
//   - args: ...any - The arguments to format.
func PrintErrorSf(err string, args ...any) {
	globalPrinter.WriteToErr([]byte(fmt.Sprintf(err, args...)))
}

// Errorf logs a formatted error message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func Errorf(format string, a ...any) {
	globalPrinter.Errorf(format, a...)
}

// Warnf logs a formatted warning message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func Warnf(format string, a ...any) {
	globalPrinter.Warnf(format, a...)
}

// Infof logs a formatted informational message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func Infof(format string, a ...any) {
	globalPrinter.Infof(format, a...)
}

// Debugf logs a formatted debug message using the global printer.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func Debugf(format string, a ...any) {
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

// WithField adds a single key-value pair to the global printer's context.
//
// Parameters:
//   - key: string - The key for the field.
//   - value: any - The value associated with the key.
//
// Returns:
//   - *Printer: A new Printer instance with the added field.
func WithField(key string, value any) *Printer {
	return globalPrinter.WithField(key, value)
}

// WithFields adds multiple key-value pairs to the global printer's context.
//
// Parameters:
//   - fields: LogFields - A map containing the fields to add.
//
// Returns:
//   - *Printer: A new Printer instance with the added fields.
func WithFields(fields LogFields) *Printer {
	return globalPrinter.WithFields(fields)
}
