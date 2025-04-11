# Printer Package

The `printer` package provides a comprehensive logging utility with color formatting for console output. It allows for various log levels and formatted output to standard output and error streams.

## Features

- Log messages with different levels: Error, Warn, Info, Debug.
- Color formatting for log messages.
- Thread-safe logging.
- Global printer instance for easy logging.

## Installation

To install the `printer` package, use `go get`:

```sh
go get github.com/cruffinoni/printer
```

## Usage

### Creating a Printer

To create a new `Printer` instance, use the `NewPrint` function:

```go
import (
    "os"
    "github.com/cruffinoni/printer"
)

printer := printer.NewPrint(printer.LevelDebug, os.Stdin, os.Stdout, os.Stderr)
```

### Logging Methods

#### Global Printer Functions

The package provides a global printer instance for convenience:

```go
printer.Printf("Hello, %s!", "world")
printer.Print("This is a standard message.")
printer.PrintError(fmt.Errorf("This is an error message."))
printer.Errorf("This is a formatted error message: %s", "error details")
printer.Warnf("This is a warning message.")
printer.Infof("This is an info message.")
printer.Debugf("This is a debug message.")
```

#### Printer Methods

For more control, you can use the `Printer` methods:

```go
printer.WriteToStd([]byte("Standard message"))
printer.WriteToError([]byte("Error message"))
printer.Errorf("Formatted error message: %s", "error details")
printer.Warnf("Warning message")
printer.Infof("Info message")
printer.Debugf("Debug message")
```

### Setting and Getting Log Level

To set the log level:

```go
printer.SetLogLevel(printer.LevelInfo)
```

To get the current log level:

```go
logLevel := printer.GetLogLevel()
fmt.Println("Current log level:", logLevel)
```

### Closing the Printer

To close the `Printer` and its associated files:

```go
err := printer.Close()
if err != nil {
    fmt.Println("Error closing printer:", err)
}
```

## Log Levels

The package defines four log levels:

- `LevelError` (0)
- `LevelWarn` (1)
- `LevelInfo` (2)
- `LevelDebug` (3)

## Color Formatting

The package supports color formatting using special tags:

- Foreground Colors: `F_<color>`
- Background Colors: `B_<color>`
- Options: `BOLD`, `FAINT`, `UNDERLINED`, `SLOWBLINK`, `RESET`
- Colors: `BLACK`, `RED`, `GREEN`, `YELLOW`, `BLUE`, `MAGENTA`, `CYAN`, `WHITE`

Example:

```go
message := "{{{-F_RED,BOLD}}}This is a bold red message{{{-RESET}}}"
printer.Print(message)
```
## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
