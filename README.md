# Printer Package

The `printer` package provides a comprehensive logging utility with color formatting for console output. It allows for various log levels and formatted output to standard output and error streams.

## Features

- Log messages with different levels: Error, Warn, Info, Debug.
- Color formatting for log messages.
- Thread-safe logging.
- Global printer instance for easy logging.
- Structured logging with fields.
- Log truncation for long messages.
- Field truncation for long fields.

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

printer := printer.NewPrint(printer.LevelDebug, printer.FlagWithDate|printer.FlagWithGoroutineID, os.Stdin, os.Stdout, os.Stderr)
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

### Structured Logging

The `Printer` struct supports structured logging with fields. You can add fields to log entries using the `WithField` and `WithFields` methods.

#### Adding a Single Field

To add a single field to a log entry, use the `WithField` method:

```go
printer.WithField("key", "value").Infof("This is an info message with a field")
```

#### Adding Multiple Fields

To add multiple fields to a log entry, use the `WithFields` method:

```go
fields := printer.LogFields{
    "key1": "value1",
    "key2": "value2",
}
printer.WithFields(fields).Infof("This is an info message with multiple fields")
```

#### Printing Out Fields

When there are fields added to a log entry, the `Printer` struct will print out the fields in the log message. The fields will be included in the log prefix.

Example:

```go
fields := printer.LogFields{
    "user": "john_doe",
    "action": "login",
}
printer.WithFields(fields).Infof("User action logged")
```

The log message will include the fields in the prefix:

```
[15:04:05.000 | INFO | user=john_doe, action=login] User action logged
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

### Setting Maximum Log Length

To set the maximum log length for truncation:

```go
printer.SetMaxLogLength(100)
```

### Setting Maximum Field Length

To set the maximum field length for truncation:

```go
printer.SetMaxFieldLength(50)
```

### Closing the Printer

To close the `Printer` and its associated files:

```go
err := printer.Close()
if err != nil {
    fmt.Println("Error closing printer:", err)
}
```

### Flags

The `Printer` struct includes flags to control the logging behavior. The available flags are:

- `FlagWithDate`: Include the current date and time in log messages.
- `FlagWithGoroutineID`: Include the goroutine ID in log messages.
- `FlagWithColor`: Enable color formatting in log messages (enabled by default).
- `FlagTruncateLogs`: Enable log truncation.
- `FlagTruncateFields`: Enable field truncation.

To create a new `Printer` instance with specific flags, use the `NewPrint` function and combine the flags using the bitwise OR operator (`|`):

```go
printer := printer.NewPrint(printer.LevelDebug, printer.FlagWithDate|printer.FlagWithGoroutineID, os.Stdout, os.Stderr)
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
printer.Print("{{{-F_RED,BOLD}}}This is a bold red message{{{-RESET}}}")
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Modules](https://blog.golang.org/using-go-modules)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Concurrency Patterns](https://blog.golang.org/pipelines)
