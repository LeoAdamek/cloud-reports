package printer

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Color represents a pen color
type Color int

type Mode int

const (
	BLACK Color = iota
	BLUE
	GREEN
	RED
	
	ModeText  Mode = 0
	ModeGraphics   = 1
	
	crlf = "\r\n"
)

// Printer defines a printer
type Printer struct {
	// Writer for the printer
	Writer io.Writer
	
	// Total width of page in steps
	Width int
	
	// Number of steps in a millimetre
	StepsPerMM int
	
	// Internal tracking of cursor location
	cur Cursor
	
	// Mode for printing
	mode Mode
	
	
	// Log
	logger *log.Logger
	
	Debug bool
}

// Cursor Position
type Cursor struct {
	X int
	Y int
	Color Color
}

func (c Cursor) Pos() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func New(w io.Writer) Printer {
	return Printer{
		cur: Cursor{
			X: 0,
			Y: 0,
			Color: BLACK,
		},
		
		logger: log.New(os.Stdout, "printer ", log.LstdFlags),
		mode: ModeText,
		
		Writer: w,
		
		Width: 480,
		StepsPerMM: 5,
	}
}

func (p *Printer) SetMode(m Mode) {
	
	chr := byte(m + 0x11)
	
	p.log(fmt.Sprintf("Setting mode to: 0x%02x", chr))
	
	p.Writer.Write([]byte{chr, '\r', '\n'})
	
	p.mode = m
}


func (p *Printer) SetOrigin() {
	p.cur.X = 0
	p.cur.Y = 0
	
	p.Writer.Write([]byte("I" + crlf))

	p.log("Set Home")
}

func (p *Printer) Home() {
	p.cur.X = 0
	p.cur.Y = 0
	
	p.Writer.Write([]byte("H" + crlf))

	p.log("Home")
}

func (p *Printer) Reset() {
	p.cur.X = 0
	p.cur.Y = 0
	p.mode = ModeText
	p.Writer.Write([]byte("A" + crlf))

	p.log("Reset Printer")
}

func (p Printer) SetColor(color Color) {
	p.cur.Color = color
	
	cmd := fmt.Sprintf("C%d", color)
	
	p.log(cmd)
	
	p.Writer.Write([]byte(cmd + crlf))
}

func (p Printer) CurrentMode() Mode {
	return p.mode
}

func (p Printer) CurrentCursor() Cursor {
	return p.cur
}




func (p *Printer) log(msg string) {
	if p.Debug {
		p.logger.Println(msg)
	}
}