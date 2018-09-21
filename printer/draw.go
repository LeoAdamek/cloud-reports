package printer

import (
	"errors"
	"fmt"
)

type TextOrientation int

const (
	TextLeftToRight TextOrientation = iota
	TextTopToBottom
	TextRightToLeft
	TextBottomToTop
)

func (p *Printer) Move(x, y int) error {
	dx := p.cur.X + x
	dy := p.cur.Y + y
	
	return p.MoveTo(dx, dy)
}

func (p *Printer) MoveTo(x, y int) error {
	return p.translate("M", x, y)
}


func (p *Printer) Draw(x, y int) error {
	dx := p.cur.X + x
	dy := p.cur.Y + y
	
	return p.DrawTo(dx, dy)
}

func (p *Printer) DrawTo(x, y int) error {
	return p.translate("D", x, y)
}


func (p *Printer) SetLineMode(l int) error {
	if l < 0 || l > 15 {
		return errors.New("invalid line type, must bet 0..15")
	}
	
	cmd := fmt.Sprintf("L%d", l)
	
	p.log(cmd)
	
	p.Writer.Write([]byte(cmd + crlf))
	
	return nil
}


func (p *Printer) GraphicalWrite(text string, size int, orientation TextOrientation) {
	
	cmd := fmt.Sprintf("Q%d", orientation)
	p.log(cmd)
	p.Writer.Write([]byte(cmd + crlf))
	
	cmd = fmt.Sprintf("S%d", size)
	p.log(cmd)
	p.Writer.Write([]byte(cmd + crlf))
	
	cmd = "P" + text
	p.log(cmd)
	p.Writer.Write([]byte(cmd + crlf))
}

func (p *Printer) translate(t string, x, y int) error {
	
	dst := Cursor{
		X: x,
		Y: y,
		Color: p.cur.Color,
	}
	
	cmd := fmt.Sprintf("%s%d,%d", t, dst.X, dst.Y)
	
	p.log(cmd)
	
	p.cur = dst
	
	_, err := p.Writer.Write([]byte(cmd + crlf))
	
	return err
}
