package printer

type Levels int

const (
	LevelError Levels = iota // Error level logging
	LevelWarn                // Warning level logging
	LevelInfo                // Info level logging
	LevelDebug               // Debug level logging
)
