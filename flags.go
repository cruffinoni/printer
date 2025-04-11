package printer

type Flags uint

const (
	FlagWithDate Flags = 1 << iota
	FlagWithGoroutineID
	FlagWithColor
)
