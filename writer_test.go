package printer

import (
	"log"
	"testing"
)

func TestNewPrint(t *testing.T) {
	p := NewPrint(LevelDebug, nil, nil, nil)
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
	p := NewPrint(LevelDebug, nil, nil, nil)
	buffer := []byte("test > {{{F_RED}}}red{{{RESET}}} < test")
	res := p.formatColor(buffer)
	log.Printf("res: %s", res)
}
