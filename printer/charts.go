package printer

import (
	"fmt"
)

type Direction int

const (
	 DirVertical Direction = iota
	 DirHorizontal
)



// Axis draws axis
func (p Printer) Axis(dir Direction, ticksEvery, nTicks int) error {
	
	length := ticksEvery * nTicks
	
	dst := p.cur
	
	if dir == DirHorizontal {
		dst.X += length
		
		p.cur.X += ticksEvery * nTicks
	} else {
		dst.Y += length
		
		
		p.cur.Y += ticksEvery * nTicks
	}
	
	cmd := fmt.Sprintf("X%d,%d,%d", dir, ticksEvery, nTicks)
	
	p.log(cmd)
	
	
	_, err := p.Writer.Write([]byte(cmd + crlf))
	
	return err
}


