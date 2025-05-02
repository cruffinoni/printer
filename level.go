package printer

// Levels represents the logging levels used in the Printer package.
type Levels int

const (
	// LevelError represents error level logging.
	LevelError Levels = iota
	// LevelWarn represents warning level logging.
	LevelWarn
	// LevelInfo represents info level logging.
	LevelInfo
	// LevelDebug represents debug level logging.
	LevelDebug
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
