package cmd

import (
	"github.com/urfave/cli"
	"math"
	"reports/printer"
	"time"
)

func testChart(c *cli.Context) {
	
	p := getPrinter(c)
	
	tc := printer.TimeSeriesChart{}
	
	tc.SampleInterval = 72 * time.Second
	tc.End = time.Now().Truncate(24*time.Hour)
	tc.Start = tc.End.AddDate(0,0,-1)
	
	sin := make([]float64, 600)
	cos := make([]float64, 600)
	
	for i := 0; i < 600; i++ {
		sin[i] = math.Sin(float64(i) * 0.01 * math.Pi) + 1
		cos[i] = math.Cos(float64(i) * 0.02 * math.Pi) + 1
	}
	
	tc.Series = [][]float64{sin, cos}
	
	p.DrawTimeSeries(tc)
}