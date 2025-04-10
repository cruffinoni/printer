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
func (l *Printer) formatColor(buffer []byte) []byte {
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
func (l *Printer) WriteToError(b []byte) {
	l.write(append([]byte("{{{-F_RED,BOLD}}}Error:{{{-RESET}}} "), b...), l.err)
}

// WriteToStd writes a raw message to the standard output stream.
func (l *Printer) WriteToStd(b []byte) {
	l.write(b, l.out)
}

// write writes a byte slice to the specified output stream with formatting and locking.
func (l *Printer) write(b []byte, out *os.File) {
	l.mx.Lock()
	defer l.mx.Unlock()
	b = l.formatColor(b)
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
func (l *Printer) WriteToStdf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	l.write(b, l.out)
}

// WriteToErrf formats a message using fmt.Sprintf and writes it to the error output stream.
func (l *Printer) WriteToErrf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	l.WriteToError(b)
}

// SetLogLevel updates the log level of the Printer.
func (l *Printer) SetLogLevel(level int) {
	l.logLevel = level
}

// GetLogLevel retrieves the current log level.
func (l *Printer) GetLogLevel() int {
	return l.logLevel
}

// formatPrefix returns a formatted log prefix with goroutine ID, timestamp, and log level.
func (l *Printer) formatPrefix(level string) string {
	return fmt.Sprintf("[%03d | %s | %s]", getGoroutineID(), time.Now().Format("15:04:05.000"), level)
}

// Errorf logs an error message if the log level permits.
func (l *Printer) Errorf(format string, a ...interface{}) {
	if l.logLevel >= LevelError {
		msg := fmt.Sprintf("{{{-F_RED,BOLD}}}"+l.formatPrefix("ERROR")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.err)
	}
}

// Warnf logs a warning message if the log level permits.
func (l *Printer) Warnf(format string, a ...interface{}) {
	if l.logLevel >= LevelWarn {
		msg := fmt.Sprintf("{{{-F_YELLOW,BOLD}}}"+l.formatPrefix("WARN")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.out)
	}
}

// Infof logs an informational message if the log level permits.
func (l *Printer) Infof(format string, a ...interface{}) {
	if l.logLevel >= LevelInfo {
		msg := fmt.Sprintf("{{{-F_BLUE,BOLD}}}"+l.formatPrefix("INFO")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.out)
	}
}

// Debugf logs a debug message if the log level permits.
func (l *Printer) Debugf(format string, a ...interface{}) {
	if l.logLevel >= LevelDebug {
		msg := fmt.Sprintf("{{{-F_CYAN,BOLD}}}"+l.formatPrefix("DEBUG")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.out)
	}
}

// Close safely closes all associated I/O streams of the Printer.
func (l *Printer) Close() error {
	if l.out != nil {
		err := l.out.Close()
		if err != nil {
			return err
		}
	}
	if l.err != nil {
		err := l.err.Close()
		if err != nil {
			return err
		}
	}
	if l.in != nil {
		err := l.in.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
