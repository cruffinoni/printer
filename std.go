package printer

import "os"

var globalPrinter = NewPrint(LevelDebug, FlagWithDate|FlagWithGoroutineID|FlagWithColor, os.Stdin, os.Stdout, os.Stderr)

func Printf(p string, a ...interface{}) {
	globalPrinter.WriteToStdf(p, a...)
}

func Print(s string) {
	globalPrinter.WriteToStd([]byte(s))
}

func PrintError(err error) {
	if err == nil {
		globalPrinter.WriteToError([]byte("<nil>"))
	} else {
		globalPrinter.WriteToError([]byte(err.Error()))
	}
}

func PrintErrorS(err string) {
	globalPrinter.WriteToError([]byte(err))
}

func PrintErrorSf(err string, args ...interface{}) {
	globalPrinter.WriteToErrf(err, args...)
}

func Errorf(format string, a ...interface{}) {
	globalPrinter.Errorf(format, a...)
}

func Warnf(format string, a ...interface{}) {
	globalPrinter.Warnf(format, a...)
}

func Infof(format string, a ...interface{}) {
	globalPrinter.Infof(format, a...)
}

func Debugf(format string, a ...interface{}) {
	globalPrinter.Debugf(format, a...)
}

func SetLogLevel(level int) {
	globalPrinter.SetLogLevel(level)
}

func GetLogLevel() int {
	return globalPrinter.GetLogLevel()
}

func WithField(key string, value interface{}) *Printer {
	return globalPrinter.WithField(key, value)
}

func WithFields(fields map[string]interface{}) *Printer {
	return globalPrinter.WithFields(fields)
}
