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

### Creating a Writer

To create a new `Writer` instance, use the `NewPrint` function:

```go
import (
    "os"
    "github.com/cruffinoni/printer"
)

writer := printer.NewPrint(printer.LevelDebug, os.Stdin, os.Stdout, os.Stderr)
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

#### Writer Methods

For more control, you can use the `Writer` methods:

```go
writer.WriteToStd([]byte("Standard message"))
writer.WriteToError([]byte("Error message"))
writer.Errorf("Formatted error message: %s", "error details")
writer.Warnf("Warning message")
writer.Infof("Info message")
writer.Debugf("Debug message")
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
