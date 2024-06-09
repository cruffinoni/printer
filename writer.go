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

type Writer struct {
	out      *os.File
	in       *os.File
	err      *os.File
	logLevel int
	mx       *sync.RWMutex
}

func NewPrint(loglevel int, in, out, err *os.File) *Writer {
	return &Writer{
		out:      out,
		in:       in,
		err:      err,
		logLevel: loglevel,
		mx:       &sync.RWMutex{},
	}
}

const (
	prefixB = "B_"
	prefixF = "F_"
)

const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

var colorFinderRegex = regexp.MustCompile(`\{{3}-?([\w,_]*)}{3}`)

func (l *Writer) formatColor(buffer []byte) []byte {
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

		output.Truncate(output.Len() - 1) // Remove the last semicolon
		output.WriteByte('m')

		buffer = bytes.ReplaceAll(buffer, i[0], output.Bytes())
		output.Reset()
	}

	buffer = append(buffer, []byte("\x1b[0m")...)
	return buffer
}

func (l *Writer) WriteToError(b []byte) {
	l.write(append([]byte("{{{-F_RED,BOLD}}}Error:{{{-RESET}}} "), b...), l.err)
}

func (l *Writer) WriteToStd(b []byte) {
	l.write(b, l.out)
}

func (l *Writer) write(b []byte, out *os.File) {
	l.mx.RLock()
	b = l.formatColor(b)
	defer l.mx.RUnlock()
	bt := []byte("\n")
	if !bytes.HasSuffix(b, bt) {
		b = append(b, bt...)
	}
	_, err := out.Write(b)
	if err != nil {
		panic(err)
	}
}

func (l *Writer) WriteToStdf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	l.write(b, l.out)
}

func (l *Writer) WriteToErrf(format string, a ...any) {
	b := []byte(fmt.Sprintf(format, a...))
	l.WriteToError(b)
}

func (l *Writer) SetLogLevel(level int) {
	l.logLevel = level
}

func (l *Writer) GetLogLevel() int {
	return l.logLevel
}

func (l *Writer) formatPrefix(level string) string {
	return fmt.Sprintf("[%03d | %s | %s]", getGoroutineID(), time.Now().Format("15:04:05.000"), level)
}

func (l *Writer) Errorf(format string, a ...interface{}) {
	if l.logLevel >= LevelError {
		msg := fmt.Sprintf("{{{-F_RED,BOLD}}}"+l.formatPrefix("ERROR")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.err)
	}
}

func (l *Writer) Warnf(format string, a ...interface{}) {
	if l.logLevel >= LevelWarn {
		msg := fmt.Sprintf("{{{-F_YELLOW,BOLD}}}"+l.formatPrefix("WARN")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.out)
	}
}

func (l *Writer) Infof(format string, a ...interface{}) {
	if l.logLevel >= LevelInfo {
		msg := fmt.Sprintf("{{{-F_BLUE,BOLD}}}"+l.formatPrefix("INFO")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.out)
	}
}

func (l *Writer) Debugf(format string, a ...interface{}) {
	if l.logLevel >= LevelDebug {
		msg := fmt.Sprintf("{{{-F_CYAN,BOLD}}}"+l.formatPrefix("DEBUG")+" {{{-RESET}}}"+format, a...)
		l.write([]byte(msg), l.out)
	}
}
