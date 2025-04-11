package printer

type Levels int

const (
	LevelError Levels = iota // Error level logging
	LevelWarn                // Warning level logging
	LevelInfo                // Info level logging
	LevelDebug               // Debug level logging
)

// String returns the string representation of the logging level.
func (l Levels) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// GetColor returns the color code for the given logging level.
func (l Levels) GetColor() string {
	switch l {
	case LevelError:
		return "F_RED,BOLD"
	case LevelWarn:
		return "F_YELLOW,BOLD"
	case LevelInfo:
		return "F_GREEN,BOLD"
	case LevelDebug:
		return "F_BLUE,BOLD"
	default:
		return "F_WHITE,BOLD"
	}
}
