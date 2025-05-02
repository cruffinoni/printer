package printer

import (
	"runtime"
	"strconv"
	"strings"
)

// getGoroutineID retrieves the ID of the current goroutine.
//
// This function uses the runtime.Stack function to obtain the stack trace of the current goroutine,
// extracts the goroutine ID from the stack trace, and returns it as a uint64 value.
//
// Returns:
//   - uint64: The ID of the current goroutine.
func getGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	s := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseUint(s, 10, 64)
	return id
}
