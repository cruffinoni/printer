package printer

// Flags represents a set of configurable options for the Printer.
//
// This type is defined as an unsigned integer and uses bitwise operations
// to enable multiple flags to be combined and checked efficiently.
type Flags uint

const (
	// FlagWithDate enables the inclusion of the current date in the output.
	FlagWithDate Flags = 1 << iota

	// FlagWithGoroutineID enables the inclusion of the current goroutine ID in the output.
	FlagWithGoroutineID

	// FlagWithColor enables colored output for better readability.
	FlagWithColor
)
