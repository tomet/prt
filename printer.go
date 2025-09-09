package prt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var (
	newLine = []byte("\n")
)

func NewPrinter(w io.Writer) *Printer {
	if c, ok := w.(io.Closer); ok {
		return &Printer{
			writer: w,
			closer: c,
		}
	}
	return &Printer{writer: w}
}

func NewBufferedPrinter(w io.Writer) *Printer {
	if c, ok := w.(io.Closer); ok {
		return &Printer{
			writer: bufio.NewWriter(w),
			closer: c,
		}
	}
	return &Printer{writer: bufio.NewWriter(w)}
}

//--------------------------------------------------------------------------------
// fmt.Stringer
//--------------------------------------------------------------------------------

// Bei einem StringPrinter wird der gebildete String geliefert.
//
// Bei allen anderen Printern wird fmt.Sprintf("Printer(%v %v)", p.writer, p.closer) geliefert.
func (p *Printer) String() string {
	switch v := p.writer.(type) {
	case *strings.Builder:
		return v.String()
	default:
		return fmt.Sprintf("Printer(%v, %v)", p.writer, p.closer)
	}
}

//--------------------------------------------------------------------------------
// Print*
//--------------------------------------------------------------------------------

func (p *Printer) Print(args ...any) {
	for _, arg := range args {
		p.WriteString(ToStr(arg))
	}
}

func (p *Printer) Println(args ...any) {
	p.Print(args...)
	p.NewLine()
}

func (p *Printer) Printf(format string, args ...any) {
	p.WriteString(fmt.Sprintf(format, args...))
}

func (p *Printer) Printfln(format string, args ...any) {
	p.WriteString(fmt.Sprintf(format, args...))
	p.NewLine()
}

//--------------------------------------------------------------------------------
// NewLine
//--------------------------------------------------------------------------------

func (p *Printer) NewLine() {
	p.Write(newLine)
}

//--------------------------------------------------------------------------------
// Write*
//--------------------------------------------------------------------------------

func (p *Printer) WriteByte(b byte) {
	p.Write([]byte{b})
}

func (p *Printer) WriteRune(r rune) {
	p.WriteString(string([]rune{r}))
}

func (p *Printer) WriteString(s string) {
	p.Write([]byte(s))
}

// Die Basis-Funktion für alle Ausgaben. Alle anderen Print*- oder Write*-Funktionen
// verwenden [*Printer.Write].
//
// Implementiert das [io.Writer] interface.
func (p *Printer) Write(bytes []byte) (n int, err error) {
	if p.err != nil || p.writer == nil {
		return 0, p.err
	}

	n, err = p.writer.Write(bytes)
	if err != nil {
		p.err = err
		p.Close()
	}

	return n, err
}

//--------------------------------------------------------------------------------
// Err
//--------------------------------------------------------------------------------

func (p *Printer) Err() error {
	return p.err
}

//--------------------------------------------------------------------------------
// Flush
//--------------------------------------------------------------------------------

// Falls der zugrundeliegende [io.Writer] ein [*bufio.Writer] (siehe [NewBufferedPrinter] und
// [NewFilePrinter]) ist, wird [*bufio.Flush] aufgerufen.
func (p *Printer) Flush() {
	if p.writer == nil {
		return
	}

	if w, ok := p.writer.(*bufio.Writer); ok {
		err := w.Flush()
		if err != nil && p.err == nil {
			p.err = err
			p.Close()
		}
	}
}

//--------------------------------------------------------------------------------
// Close
//--------------------------------------------------------------------------------

// Falls der zugrundeliegende [io.Writer] auch ein [io.Closer] ist,
// wird dieser nun geschlossen.
//
// Falls der zugrundeliegende [io.Writer] auch ein [*bufio.Writer] ist,
// wird zuvor noch [*bufio.Writer.Flush] aufgerufen.
func (p *Printer) Close() error {
	if p.writer != nil && p.err == nil {
		if w, ok := p.writer.(*bufio.Writer); ok {
			err := w.Flush()
			if err != nil && p.err == nil {
				p.err = err
			}
		}
	}

	p.writer = nil

	if p.closer != nil {
		err := p.closer.Close()
		p.closer = nil
		if err != nil && p.err == nil {
			p.err = err
		}
	}

	return p.err
}
