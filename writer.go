package printer

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Printer provides structured output to various I/O streams with support for
// log levels, colored output, and concurrency-safe operations.
type Printer struct {
	out      *os.File   // Output stream for standard messages
	in       *os.File   // Input stream, if applicable
	err      *os.File   // Output stream for error messages
	logLevel int        // Current logging level
	mx       sync.Mutex // Mutex for synchronized writes
}

// NewPrint creates a new Printer instance with specified log level and I/O streams.
//
// Parameters:
//   - loglevel: int - The initial logging level.
//   - in: *os.File - The input stream (can be nil if not applicable).
//   - out: *os.File - The output stream for standard messages.
//   - err: *os.File - The output stream for error messages.
//
// Returns:
//   - *Printer: A new Printer instance.
func NewPrint(loglevel int, in, out, err *os.File) *Printer {
	return &Printer{
		out:      out,
		in:       in,
		err:      err,
		logLevel: loglevel,
		mx:       sync.Mutex{},
	}
}

const (
	prefixB = "B_" // Prefix for background colors
	prefixF = "F_" // Prefix for foreground colors
)

const (
	LevelError = iota // Error level logging
	LevelWarn         // Warning level logging
	LevelInfo         // Info level logging
	LevelDebug        // Debug level logging
)

// bufferPool provides reusable byte buffers to reduce memory allocations.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// colorFinderRegex matches color formatting placeholders in the log strings.
var colorFinderRegex = regexp.MustCompile(`\{{3}-?([\w,_]*)}{3}`)

// formatColor replaces color formatting tokens in the buffer with ANSI codes.
//
// Parameters:
//   - buffer: []byte - The input buffer containing color formatting tokens.
//
// Returns:
//   - []byte: The buffer with color formatting tokens replaced by ANSI codes.
func (p *Printer) formatColor(buffer []byte) []byte {
	f := colorFinderRegex.FindAllSubmatch(buffer, -1)
	if f == nil {
		return buffer
	}

	output := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(output)
	output.Reset()

	for _, i := range f {
		output.WriteString("\x1b[")

		composed := bytes.Split(i[1], []byte(","))
		for _, c := range composed {
			if bytes.HasPrefix(c, []byte(prefixB)) {
				color := bytes.TrimPrefix(c, []byte(prefixB))
				if col, ok := colorValues[strings.ToLower(string(color))]; ok {
					_, _ = fmt.Fprintf(output, "%d;", col+BackgroundBlack)
				} else {
					_, _ = fmt.Fprintf(output, "%%B_COLOR_NOT_FOUND%%%s%%", c)
				}
			} else if bytes.HasPrefix(c, []byte(prefixF)) {
				color := bytes.TrimPrefix(c, []byte(prefixF))
				if col, ok := colorValues[strings.ToLower(string(color))]; ok {
					_, _ = fmt.Fprintf(output, "%d;", col+ForegroundBlack)
				} else {
					_, _ = fmt.Fprintf(output, "%%F_COLOR_NOT_FOUND%%%s%%", c)
				}
			} else {
				if opt, ok := colorOptions[strings.ToLower(string(c))]; ok {
					_, _ = fmt.Fprintf(output, "%d;", opt)
				} else {
					_, _ = fmt.Fprintf(output, "%%NOT_FOUND%%%s%%", c)
				}
			}
		}

		output.Truncate(output.Len() - 1)
		output.WriteByte('m')

		buffer = bytes.ReplaceAll(buffer, i[0], output.Bytes())
		output.Reset()
	}

	buffer = append(buffer, []byte("\x1b[0m")...)
	return buffer
}

// WriteToError writes a formatted error message to the error output stream.
//
// Parameters:
//   - b: []byte - The error message to write.
func (p *Printer) WriteToError(b []byte) {
	p.write(append([]byte("{{{-F_RED,BOLD}}}Error:{{{-RESET}}} "), b...), p.err)
}

// WriteToStd writes a raw message to the standard output stream.
//
// Parameters:
//   - b: []byte - The message to write.
func (p *Printer) WriteToStd(b []byte) {
	p.write(b, p.out)
}

// Write writes a byte slice to the standard output stream.
//
// Parameters:
//   - buffer: []byte - The data to write.
//
// Returns:
//   - n: int - The number of bytes written.
//   - err: error - Any error encountered during the write operation.
func (p *Printer) Write(buffer []byte) (n int, err error) {
	if p.out != nil {
		n, err = p.out.Write(buffer)
	}
	return
}

// write writes a byte slice to the specified output stream with formatting and locking.
//
// Parameters:
//   - b: []byte - The data to write.
//   - out: *os.File - The output stream to write to.
func (p *Printer) write(b []byte, out *os.File) {
	p.mx.Lock()
	defer p.mx.Unlock()
	b = p.formatColor(b)
	bt := []byte("\n")
	if !bytes.HasSuffix(b, bt) {
		b = append(b, bt...)
	}
	_, err := out.Write(b)
	if err != nil {
		panic(err)
	}
}

// WriteToStdf formats a message using fmt.Sprintf and writes it to the standard output stream.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func (p *Printer) WriteToStdf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	p.write(b, p.out)
}

// WriteToErrf formats a message using fmt.Sprintf and writes it to the error output stream.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func (p *Printer) WriteToErrf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	p.WriteToError(b)
}

// SetLogLevel updates the log level of the Printer.
//
// Parameters:
//   - level: int - The new log level to set.
func (p *Printer) SetLogLevel(level int) {
	p.logLevel = level
}

// GetLogLevel retrieves the current log level.
//
// Returns:
//   - int: The current log level.
func (p *Printer) GetLogLevel() int {
	return p.logLevel
}

// formatPrefix returns a formatted log prefix with goroutine ID, timestamp, and log level.
//
// Parameters:
//   - level: string - The log level as a string.
//
// Returns:
//   - string: The formatted log prefix.
func (p *Printer) formatPrefix(level string) string {
	return fmt.Sprintf("[%03d | %s | %s]", getGoroutineID(), time.Now().Format("15:04:05.000"), level)
}

// Errorf logs an error message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func (p *Printer) Errorf(format string, a ...interface{}) {
	if p.logLevel >= LevelError {
		msg := fmt.Sprintf("{{{-F_RED,BOLD}}}"+p.formatPrefix("ERROR")+" {{{-RESET}}}"+format, a...)
		p.write([]byte(msg), p.err)
	}
}

// Warnf logs a warning message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func (p *Printer) Warnf(format string, a ...interface{}) {
	if p.logLevel >= LevelWarn {
		msg := fmt.Sprintf("{{{-F_YELLOW,BOLD}}}"+p.formatPrefix("WARN")+" {{{-RESET}}}"+format, a...)
		p.write([]byte(msg), p.out)
	}
}

// Infof logs an informational message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func (p *Printer) Infof(format string, a ...interface{}) {
	if p.logLevel >= LevelInfo {
		msg := fmt.Sprintf("{{{-F_BLUE,BOLD}}}"+p.formatPrefix("INFO")+" {{{-RESET}}}"+format, a...)
		p.write([]byte(msg), p.out)
	}
}

// Debugf logs a debug message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...interface{} - The arguments to format.
func (p *Printer) Debugf(format string, a ...interface{}) {
	if p.logLevel >= LevelDebug {
		msg := fmt.Sprintf("{{{-F_CYAN,BOLD}}}"+p.formatPrefix("DEBUG")+" {{{-RESET}}}"+format, a...)
		p.write([]byte(msg), p.out)
	}
}

// Close safely closes all associated I/O streams of the Printer.
//
// Returns:
//   - error: Any error encountered during the close operation, or nil if successful.
func (p *Printer) Close() error {
	if p.out != nil {
		err := p.out.Close()
		if err != nil {
			return err
		}
	}
	if p.err != nil {
		err := p.err.Close()
		if err != nil {
			return err
		}
	}
	if p.in != nil {
		err := p.in.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
