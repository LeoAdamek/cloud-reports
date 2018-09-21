package printer

import (
	"fmt"
	"math"
	"time"
)

// Specialist functions for drawing charts of timeseries data


// TimeSeriesChart describes a chart containing multiple series of _aligned_ TimeSeries data
type TimeSeriesChart struct {
	Start time.Time
	End time.Time
	SampleInterval time.Duration
	Series [][]float64
	Title string
	Unit string
}


// Min gets the minimum value in across all series
func (t TimeSeriesChart) Min() float64 {
	min := math.MaxFloat64
	
	for _, r := range t.Series {
		for _, v := range r {
			if v < min {
				min = v
			}
		}
	}
	
	return min
}

// Max gets the maximum value across all series
func (t TimeSeriesChart) Max() float64 {
	max := -math.MaxFloat64
	
	for _, r := range t.Series {
		for _, v := range r {
			if v > max {
				max = v
			}
		}
	}
	
	return max
}

// Duration gets the total duration of the chart
func (t TimeSeriesChart) Duration() time.Duration{
	return t.End.Sub(t.Start)
}


func (p Printer) DrawTimeSeries(t TimeSeriesChart) {

	p.SetMode(ModeGraphics)
	p.Reset()
	p.SetMode(ModeGraphics)
	
	p.SetOrigin()
	
	p.drawTimeSeriesAxes(t)
	
	totalRange := t.Max() - t.Min() // Chart range
	ups := totalRange / 400.0 // Spread the range over the 480 horizontal steps
	sps := int(math.Ceil(600.0 / float64(len(t.Series[0]))))
	
	p.log(fmt.Sprintf("Steps per Sample (X): %d", sps))
	p.log(fmt.Sprintf("Units per step (Y): %.3f", ups))

	for i, s := range t.Series {
		p.SetColor( Color((i+1)%4) )

		p.SetLineMode( (i/4)%15 )

		p.MoveTo(400 - int(math.Floor(s[0]/ups)), -600)

		for q, v := range s {
			//p.log(fmt.Sprintf("s[%d][%d] = %.3f", i, q, v))
			x := 420 - int(math.Floor(v / ups))
			y := (q*sps) - 600
			
			p.DrawTo(x,y)
		}
	}
	
	if t.Title != "" {
		p.MoveTo(60, -400)
		p.SetColor(BLUE)
		p.GraphicalWrite(t.Title, 5, TextBottomToTop)
	}
	
}

func (p Printer) drawTimeSeriesAxes(t TimeSeriesChart) {
	
	origin := p.CurrentCursor()
	origin.Y -= 25*24
	
	p.Move(-10, -25)
	// go left
	p.MoveTo(0, origin.Y)
	
	// Draw Y axis (horizontally)
	p.Axis(DirHorizontal, 20, 20) // 400 steps
	
	// Draw X axis
	p.MoveTo(400, origin.Y)
	
	p.Axis(DirVertical, 25, 24)  // 12cm (600 steps)
	
	p.MoveTo(400, origin.Y)

	// Draw Labels (Y axis)
	dy := (t.Max() - t.Min()) / 400

	for i := 0; i < 20; i++ {
		p.MoveTo(400-(i*20), origin.Y-100)
		v := dy*(20.0*float64(i))
		p.GraphicalWrite(fmt.Sprintf("%6.3f%s", v, t.Unit), 1, TextBottomToTop)
	}
	
	// Draw Labels (X axis)
	ti := t.Duration() / 24
	
	p.log("Start: " + t.Start.Format(time.RFC1123Z))
	p.log("End: " + t.End.Format(time.RFC1123Z))
	p.log(fmt.Sprint("Interval:",ti))
	
	for i := 0; i < 25; i++ {
		v := t.Start.Add(ti * time.Duration(i))
		
		p.log("Time: " + v.Format(time.RFC3339))
		p.MoveTo(410, (origin.Y + 25*i)-4)
		p.GraphicalWrite(v.Format("15:04"), 1, TextLeftToRight)
	}
	
}