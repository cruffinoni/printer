package printer

// Reset resets all attributes to their default values.
const (
	Reset = iota
	// Bold sets the text to bold.
	Bold
	// Faint sets the text to faint (dim).
	Faint
	// Underlined underlines the text.
	Underlined = 4
	// SlowBlink sets the text to blink slowly.
	SlowBlink  = 5
)

// Foreground color constants.
const (
	ForegroundBlack = iota + 30
	ForegroundRed
	ForegroundGreen
	ForegroundYellow
	ForegroundBlue
	ForegroundMagenta
	ForegroundCyan
	ForegroundWhite
)

// Background color constants.
const (
	BackgroundBlack = iota + 40
	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundMagenta
	BackgroundCyan
	BackgroundWhite
)

// colorValues maps color names to their corresponding values.
var (
	colorValues = map[string]int{
		"black":   0,
		"red":     1,
		"green":   2,
		"yellow":  3,
		"blue":    4,
		"magenta": 5,
		"cyan":    6,
		"white":   7,
	}
	// colorOptions maps color options to their corresponding values.
	colorOptions = map[string]int{
		"reset":      Reset,
		"bold":       Bold,
		"faint":      Faint,
		"underlined": Underlined,
		"slowBlink":  SlowBlink,
	}
)
