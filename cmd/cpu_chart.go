package cmd

import (
	"github.com/urfave/cli"
	"log"
	"reports/printer"
	"reports/util"
	"strings"
	"time"
)

func cpuChart(c *cli.Context) {
	


	p := getPrinter(c)
	
	chart := printer.TimeSeriesChart{
		End: time.Now().Truncate(time.Minute),
		Title: "CPU Usage",
		Unit: "%",
	}
	
	chart.Start = chart.End.Add(-c.Duration("d"))
	chart.SampleInterval = c.Duration("d") / 600


	data, err := getCpuChartData(chart.Start, chart.End, chart.SampleInterval, c.StringSlice("group"))

	if err != nil {
		log.Fatalln("Error getting chart data:", err)
	}

	chart.Series = data

	p.DrawTimeSeries(chart)
	
}

func getCpuChartData(start, end time.Time, interval time.Duration, groupBy []string) ([][]float64, error) {

	if len(groupBy) == 0 {
		groupBy = []string{"instance"}
	}

	queryStr := "sum(rate(node_cpu_seconds_total{mode!=\"idle\"}[30m])) by (" + strings.Join(groupBy, ",") + ")"


	return util.GetChartData(queryStr, start, end)

}