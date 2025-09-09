package prt

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Printer struct {
	writer io.Writer
	closer io.Closer
	err    error
}

func NewStringPrinter() *Printer {
	return NewPrinter(new(strings.Builder))
}

func NewFilePrinter(filename string) (*Printer, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return NewBufferedPrinter(f), nil
}

// Wandelt einen beliebigen Wert in einen String.
//
// Es wird einfach fmt.Sprintf("%v") verwendet.
//
// Bei einer rune wird statt der Zahl das Zeichen selbst ausgegeben!
//
// Bei einem Rune-Slice ([]rune) wird das Slice in einen String umgewandelt.
func ToStr(v any) string {
	switch v2 := v.(type) {
	case rune:
		return fmt.Sprintf("%c", v2)
	case []rune:
		return string(v2)
	default:
		return fmt.Sprintf("%v", v)
	}
}
