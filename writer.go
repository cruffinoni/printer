package printer

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// LogFields represents a map of log fields for structured logging.
type LogFields map[string]any

// Printer provides structured output to various I/O streams with support for
// log levels, colored output, and concurrency-safe operations.
type Printer struct {
	out            io.WriteCloser
	err            io.WriteCloser
	logLevel       Levels
	flags          Flags
	mx             sync.Mutex
	fields         LogFields
	maxLogLength   int
	maxFieldLength int
}

// NewPrinter creates a new Printer instance with specified log level and I/O streams.
//
// Parameters:
//   - loglevel: int - The initial logging level.
//   - out: io.WriteCloser - The output stream for standard messages.
//   - err: io.WriteCloser - The output stream for error messages.
//
// Returns:
//   - *Printer: A new Printer instance.
func NewPrinter(loglevel Levels, flags Flags, out, err io.WriteCloser) *Printer {
	p := &Printer{
		out:      out,
		err:      err,
		logLevel: loglevel,
		flags:    flags | FlagPanicOnError,
		mx:       sync.Mutex{},
		fields:   make(LogFields),
	}

	if flags&FlagTruncateLogs != 0 {
		p.maxLogLength = DefaultMaxLogLength
	}
	if flags&FlagTruncateFields != 0 {
		p.maxFieldLength = DefaultMaxFieldLength
	}

	return p
}

const (
	prefixB = "B_" // Prefix for background colors
	prefixF = "F_" // Prefix for foreground colors
)

// bufferPool provides reusable byte buffers to reduce memory allocations.
var bufferPool = sync.Pool{
	New: func() any {
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
	if p.flags&FlagWithColor == 0 {
		return buffer
	}
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

// truncateLog truncates the log message if it exceeds the maximum log length.
//
// Parameters:
//   - b: []byte - The log message to truncate.
//
// Returns:
//   - []byte: The truncated log message.
func (p *Printer) truncateLog(s string) string {
	if p.flags&FlagTruncateLogs != 0 && p.maxLogLength > 0 && len(s) > p.maxLogLength {
		return s[:p.maxLogLength]
	}
	return s
}

// truncateField truncates the field if it exceeds the maximum field length.
//
// Parameters:
//   - field: string - The field to truncate.
//
// Returns:
//   - string: The truncated field.
func (p *Printer) truncateField(field string) string {
	if p.flags&FlagTruncateFields != 0 && p.maxFieldLength > 0 && len(field) > p.maxFieldLength {
		return field[:p.maxFieldLength]
	}
	return field
}

// writeTo writes a byte slice to the specified output stream with formatting and locking.
//
// Parameters:
//   - b: []byte - The data to write.
//   - writer: io.Writer - The output stream to write to.
//
// Returns:
//   - int: Number of bytes written.
//   - error: Error encountered during write, if any.
func (p *Printer) writeTo(b []byte, writer io.Writer) (int, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	b = p.formatColor(b)
	if p.flags&FlagWithoutNewLine == 0 {
		bt := []byte("\n")
		if !bytes.HasSuffix(b, bt) {
			b = append(b, bt...)
		}
	}
	return writer.Write(b)
}

// WriteToStd writes a raw message to the standard output stream.
//
// Parameters:
//   - b: []byte - The message to write.
func (p *Printer) WriteToStd(b []byte) {
	_, err := p.writeTo(b, p.out)
	if err != nil {
		panic(err)
	}
}

// Printf writes a formatted string to the standard output stream.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func (p *Printer) Printf(format string, a ...any) {
	p.WriteToStd([]byte(fmt.Sprintf(format, a...)))
}

// Print writes a raw string to the standard output stream.
//
// Parameters:
//   - s: string - The string to write.
func (p *Printer) Print(s string) {
	p.WriteToStd([]byte(s))
}

// WriteToErr writes a raw message to the error output stream.
//
// Parameters:
//   - b: []byte - The message to write.
func (p *Printer) WriteToErr(b []byte) {
	_, err := p.writeTo(b, p.err)
	if err != nil {
		panic(err)
	}
}

// Write writes a byte slice to the standard output stream.
//
// Parameters:
//   - buffer: []byte - The data to write.
//
// Returns:
//   - int: Number of bytes written.
//   - error: Error encountered during write.
func (p *Printer) Write(buffer []byte) (n int, err error) {
	return p.writeTo(buffer, p.out)
}

// SetLogLevel updates the log level of the Printer.
//
// Parameters:
//   - level: int - The new log level to set.
func (p *Printer) SetLogLevel(level Levels) {
	p.logLevel = level
}

// GetLogLevel retrieves the current log level.
//
// Returns:
//   - int: The current log level.
func (p *Printer) GetLogLevel() Levels {
	return p.logLevel
}

// SetMaxLogLength sets the maximum log length for truncation.
//
// Parameters:
//   - length: int - The maximum log length to set.
func (p *Printer) SetMaxLogLength(length int) {
	p.flags |= FlagTruncateLogs
	p.maxLogLength = length
}

// SetMaxFieldLength sets the maximum field length for truncation.
//
// Parameters:
//   - length: int - The maximum field length to set.
func (p *Printer) SetMaxFieldLength(length int) {
	p.flags |= FlagTruncateFields
	p.maxFieldLength = length
}

// formatPrefix returns a formatted log prefix with goroutine ID, timestamp, log level, and fields.
//
// Parameters:
//   - level: string - The log level as a string.
//
// Returns:
//   - string: The formatted log prefix.
func (p *Printer) formatPrefix(level Levels) string {
	content := make([]string, 0, 3)
	if p.flags&FlagWithGoroutineID != 0 {
		content = append(content, fmt.Sprintf("%03d", getGoroutineID()))
	}
	if p.flags&FlagWithDate != 0 {
		content = append(content, time.Now().Format("15:04:05.000"))
	}
	content = append(content, level.String())
	if len(p.fields) > 0 {
		fieldStrings := make([]string, 0, len(p.fields))
		for k, v := range p.fields {
			fieldStr := fmt.Sprintf("%s=\"%v\"", k, v)
			if str, ok := v.(string); ok {
				fieldStr = fmt.Sprintf("%s=%q", k, str)
			}
			fieldStr = p.truncateField(fieldStr)
			fieldStrings = append(fieldStrings, fieldStr)
		}
		sort.Strings(fieldStrings)
		content = append(content, strings.Join(fieldStrings, ", "))
	}
	if p.flags&FlagWithColor != 0 {
		return fmt.Sprintf("{{{%s}}}[%s]{{{-RESET}}} ", level.GetColor(), strings.Join(content, " | "))
	}
	return fmt.Sprintf("[%s] ", strings.Join(content, " | "))
}

// Errorf logs an error message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func (p *Printer) Errorf(format string, a ...any) {
	if p.logLevel >= LevelError {
		msg := fmt.Sprintf(p.formatPrefix(LevelError)+format, a...)
		p.WriteToErr([]byte(msg))
	}
}

// Warnf logs a warning message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func (p *Printer) Warnf(format string, a ...any) {
	if p.logLevel >= LevelWarn {
		msg := p.truncateLog(fmt.Sprintf(format, a...))
		p.WriteToStd([]byte(p.formatPrefix(LevelWarn) + msg))
	}
}

// Infof logs an informational message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func (p *Printer) Infof(format string, a ...any) {
	if p.logLevel >= LevelInfo {
		msg := p.truncateLog(fmt.Sprintf(format, a...))
		p.WriteToStd([]byte(p.formatPrefix(LevelInfo) + msg))
	}
}

// Debugf logs a debug message if the log level permits.
//
// Parameters:
//   - format: string - The format string.
//   - a: ...any - The arguments to format.
func (p *Printer) Debugf(format string, a ...any) {
	if p.logLevel >= LevelDebug {
		msg := p.truncateLog(fmt.Sprintf(format, a...))
		p.WriteToStd([]byte(p.formatPrefix(LevelDebug) + msg))
	}
}

// Close safely closes all associated I/O streams of the Printer.
//
// This method iterates over all associated I/O streams (`out` and `err`) and attempts to close them.
// If any stream fails to close, the method returns the encountered error. If all streams are closed
// successfully, it returns nil.
//
// Returns:
//   - error: An error encountered during the close operation, or nil if all streams are closed successfully.
func (p *Printer) Close() error {
	if p.out != nil {
		if err := p.out.Close(); err != nil {
			return err
		}
		p.out = nil
	}
	if p.err != nil {
		if err := p.err.Close(); err != nil {
			return err
		}
		p.err = nil
	}
	return nil
}

// Copy creates a new Printer instance with the same configuration as the current one.
//
// Returns:
//   - *Printer: A new Printer instance with the same configuration.
func (p *Printer) Copy() *Printer {
	cpyPrinter := &Printer{
		out:            p.out,
		err:            p.err,
		logLevel:       p.logLevel,
		flags:          p.flags,
		mx:             sync.Mutex{},
		fields:         make(LogFields),
		maxLogLength:   p.maxLogLength,
		maxFieldLength: p.maxFieldLength,
	}
	maps.Copy(cpyPrinter.fields, p.fields)
	return cpyPrinter
}

// WithoutColor creates a new Printer instance with color output disabled.
//
// Returns:
//   - *Printer: A new Printer instance with the color flag disabled.
func (p *Printer) WithoutColor() *Printer {
	newPrinter := p.Copy()
	newPrinter.flags &^= FlagWithColor
	return newPrinter
}

// WithField creates a new Printer instance with an additional single field.
//
// This method performs a deep copy of the current Printer instance and adds
// the specified key-value pair to the fields map of the new instance.
//
// Parameters:
//   - key: string - The key for the new field.
//   - value: any - The value for the new field.
//
// Returns:
//   - *Printer: A new Printer instance with the added field.
func (p *Printer) WithField(key string, value any) *Printer {
	newPrinter := p.Copy()
	newPrinter.fields[key] = value
	return newPrinter
}

// WithFields creates a new Printer instance with additional fields.
//
// This method performs a deep copy of the current Printer instance and adds
// the specified key-value pairs to the fields map of the new instance.
//
// Parameters:
//   - fields: LogFields - A map of key-value pairs to add to the fields.
//
// Returns:
//   - *Printer: A new Printer instance with the added fields.
func (p *Printer) WithFields(fields LogFields) *Printer {
	newPrinter := p.Copy()
	for key, value := range fields {
		newPrinter.fields[key] = value
	}
	return newPrinter
}

// WithoutNewLine creates a new Printer instance with the newline flag disabled.
//
// This method performs a deep copy of the current Printer instance and sets
// the WithoutNewLine flag in the new instance.
//
// Returns:
//   - *Printer: A new Printer instance with the WithoutNewLine flag set.
func (p *Printer) WithoutNewLine() *Printer {
	newPrinter := p.Copy()
	newPrinter.flags |= FlagWithoutNewLine
	return newPrinter
}
