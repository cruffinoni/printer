package printer

// Flags represents a set of configurable options for the Printer.
//
// This type is defined as an unsigned integer and uses bitwise operations
// to enable multiple flags to be combined and checked efficiently.
type Flags uint

const (
	// WithNoFlags is the default value for Flags, indicating no special options are set.
	WithNoFlags Flags = 0

	// FlagWithDate enables the inclusion of the current date in the output.
	FlagWithDate Flags = 1 << iota

	// FlagWithGoroutineID enables the inclusion of the current goroutine ID in the output.
	FlagWithGoroutineID

	// FlagWithColor enables colored output for better readability.
	FlagWithColor

	// FlagPanicOnError enables panic behavior on error conditions.
	FlagPanicOnError

	// FlagWithoutNewLine disables the automatic newline at the end of the output.
	FlagWithoutNewLine

	// FlagTruncateLogs enables truncation of log messages to a specified length.
	FlagTruncateLogs

	// FlagTruncateFields enables truncation of field values to a specified length.
	FlagTruncateFields
)

// Default values for log and field truncation
const (
	DefaultMaxLogLength   = 100
	DefaultMaxFieldLength = 50
)
