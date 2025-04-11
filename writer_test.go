package printer

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// dummyWC is a simple WriteCloser that wraps a bytes.Buffer. It also
// tracks whether Close() was called.
type dummyWC struct {
	buf    bytes.Buffer
	closed bool
	mu     sync.Mutex
}

func (d *dummyWC) Write(p []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return 0, errors.New("write on closed writer")
	}
	return d.buf.Write(p)
}

func (d *dummyWC) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.closed = true
	return nil
}

func (d *dummyWC) String() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.buf.String()
}

// --- Unit Tests ---

func TestPrinter(t *testing.T) {
	// Create two dummy WriteClosers for standard and error streams.
	stdOut := &dummyWC{}
	errOut := &dummyWC{}

	// Create a Printer with initial log level LevelDebug and flags that enable
	// colored output and timestamp.
	p := NewPrint(LevelDebug, FlagWithColor|FlagWithDate, stdOut, errOut)

	// Group sub-tests in a map for clarity.
	tests := map[string]func(t *testing.T){
		// Test WriteToStd: It should write the given message to the out stream,
		// appending a newline if needed.
		"WriteToStd": func(t *testing.T) {
			// We pass a message that does not end with a newline.
			msg := []byte("Test standard output")
			// Clear the stdout dummy before testing.
			stdOut.buf.Reset()
			p.WriteToStd(msg)
			// Since the message is formatted (with prefix and ANSI codes) we
			// check that our message is contained somewhere in the output and that
			// a newline is present at the end.
			outStr := stdOut.String()
			assert.Contains(t, outStr, "Test standard output")
			assert.True(t, bytes.HasSuffix([]byte(outStr), []byte("\n")))
		},

		// Test WriteToErr: similar to WriteToStd but writing to the error stream.
		"WriteToErr": func(t *testing.T) {
			msg := []byte("Error occurred")
			errOut.buf.Reset()
			p.WriteToErr(msg)
			outStr := errOut.String()
			assert.Contains(t, outStr, "Error occurred")
			assert.True(t, bytes.HasSuffix([]byte(outStr), []byte("\n")))
		},

		// Test Write method, which writes to the standard output stream.
		"Write": func(t *testing.T) {
			message := []byte("Direct write")
			stdOut.buf.Reset()
			n, err := p.Write(message)
			assert.NoError(t, err)
			// n should be at least the length of the message plus a newline.
			assert.GreaterOrEqual(t, n, len(message))
			assert.Contains(t, stdOut.String(), "Direct write")
		},

		// Test GetLogLevel and SetLogLevel.
		"SetGetLogLevel": func(t *testing.T) {
			// Set log level to LevelWarn.
			p.SetLogLevel(LevelWarn)
			assert.Equal(t, LevelWarn, p.GetLogLevel())

			// Reset back to LevelDebug for subsequent tests.
			p.SetLogLevel(LevelDebug)
			assert.Equal(t, LevelDebug, p.GetLogLevel())
		},

		// Test DisableColor: the returned Printer should have the FlagWithColor bit disabled.
		"DisableColor": func(t *testing.T) {
			// p already has FlagWithColor enabled. Create a new printer instance with color disabled.
			pNoColor := p.DisableColor()
			// As our dummy flags are just bit masks, ensure that the FlagWithColor bit is not set.
			assert.Equal(t, p.flags&FlagWithColor, FlagWithColor) // original has color enabled
			assert.Equal(t, pNoColor.flags&FlagWithColor, Flags(0))
		},

		// Test WithField: Verify that a new field is added without modifying the original.
		"WithField": func(t *testing.T) {
			originalFields := len(p.fields)
			newP := p.WithField("key", "value")
			// Original should remain unchanged.
			assert.Equal(t, originalFields, len(p.fields))
			// New printer should have one additional field.
			assert.Equal(t, originalFields+1, len(newP.fields))
			assert.Equal(t, "value", newP.fields["key"])
		},

		// Test WithFields: Verify that multiple fields are added without modifying the original.
		"WithFields": func(t *testing.T) {
			originalFields := len(p.fields)
			newFields := LogFields{
				"alpha": 1,
				"beta":  "two",
				"gamma": true,
			}
			newP := p.WithFields(newFields)
			// Original remains unchanged.
			assert.Equal(t, originalFields, len(p.fields))
			// New printer should have additional fields.
			assert.Equal(t, originalFields+len(newFields), len(newP.fields))
			for k, v := range newFields {
				assert.Equal(t, v, newP.fields[k])
			}
		},

		// Test Close: both out and err writers are closed and set to nil.
		"Close": func(t *testing.T) {
			// Create new dummy streams
			dOut := &dummyWC{}
			dErr := &dummyWC{}
			// Create a new Printer with these streams.
			p2 := NewPrint(LevelInfo, 0, dOut, dErr)
			// Call Close.
			err := p2.Close()
			assert.NoError(t, err)
			assert.Nil(t, p2.out)
			assert.Nil(t, p2.err)

			// Subsequent call to Close should succeed (or be a no-op).
			err = p2.Close()
			assert.NoError(t, err)
			// Also, attempt to write should panic if desired; we can catch that.
			assert.Panics(t, func() {
				p2.WriteToStd([]byte("should panic"))
			})
		},

		// Test log level output functions: Infof, Warnf, Errorf and Debugf.
		"LogMethods": func(t *testing.T) {
			// Reset output buffers.
			stdOut.buf.Reset()
			errOut.buf.Reset()
			// We set level to LevelDebug so all logging methods should run.
			p.SetLogLevel(LevelDebug)
			// Call each method.
			p.Debugf("Debug message")
			p.Infof("Info message")
			p.Warnf("Warn message")
			p.Errorf("Error message")
			// Collect outputs.
			outStr := stdOut.String()
			errStr := errOut.String()

			// Debug, Info, Warn should write to standard out.
			assert.Contains(t, outStr, "Debug message")
			assert.Contains(t, outStr, "Info message")
			assert.Contains(t, outStr, "WARN")
			// Error should write to error out.
			assert.Contains(t, errStr, "Error message")

			// Verify that message prefixes include a timestamp if FlagWithDate is set.
			// Because the timestamp is generated at runtime, we check for the format "15:04:05"
			now := time.Now().Format("15:04")
			assert.Contains(t, outStr, now)
		},
	}

	// Run each sub-test.
	for name, testFunc := range tests {
		t.Run(name, testFunc)
	}
}
