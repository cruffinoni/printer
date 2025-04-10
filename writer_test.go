package printer

import (
	"log"
	"os"
	"testing"
)

func TestNewPrint(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithDate|FlagWithGoroutineID, nil, nil, nil)
	if p == nil {
		t.Error("NewPrint() returned nil")
	}
}

func TestSetLogLevel(t *testing.T) {
	SetLogLevel(LevelDebug)
	if GetLogLevel() != LevelDebug {
		t.Error("SetLogLevel() failed")
	}
}

func TestGetLogLevel(t *testing.T) {
	SetLogLevel(LevelDebug)
	if GetLogLevel() != LevelDebug {
		t.Error("GetLogLevel() failed")
	}
}

func TestFormatColor(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithColor, nil, nil, nil)
	buffer := []byte("test > {{{F_RED}}}red{{{RESET}}} < test")
	res := p.formatColor(buffer)
	log.Printf("res: %s", res)
}

func TestClose(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithDate|FlagWithGoroutineID, nil, nil, nil)
	err := p.Close()
	if err != nil {
		t.Error("Close() failed")
	}
}

func TestFlagWithDate(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithDate, os.Stdin, os.Stdout, os.Stderr)
	p.Infof("This is a test message with date")
}

func TestFlagWithGoroutineID(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithGoroutineID, os.Stdin, os.Stdout, os.Stderr)
	p.Infof("This is a test message with goroutine ID")
}

func TestFlagWithColor(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithColor, os.Stdin, os.Stdout, os.Stderr)
	p.Infof("This is a test message with color")
}

func TestCombinedFlags(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithDate|FlagWithGoroutineID|FlagWithColor, os.Stdin, os.Stdout, os.Stderr)
	p.Infof("This is a test message with combined flags")
}

func TestWithField(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithDate|FlagWithGoroutineID|FlagWithColor, os.Stdin, os.Stdout, os.Stderr)
	p = p.WithField("key", "value")
	if p.fields["key"] != "value" {
		t.Error("WithField() failed to add the field")
	}
}

func TestWithFields(t *testing.T) {
	p := NewPrint(LevelDebug, FlagWithDate|FlagWithGoroutineID|FlagWithColor, os.Stdin, os.Stdout, os.Stderr)
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	p = p.WithFields(fields)
	if p.fields["key1"] != "value1" || p.fields["key2"] != "value2" {
		t.Error("WithFields() failed to add the fields")
	}
}
