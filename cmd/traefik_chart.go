package cmd

import (
	"github.com/urfave/cli"
	"log"
	"reports/printer"
	"reports/util"
	"time"
)

func traefikChart(c *cli.Context) {

	p := getPrinter(c)

	chart := printer.TimeSeriesChart{
		End: time.Now().Truncate(time.Minute),
		Title: "HTTP Requests",
		Unit: "r/s",
	}

	chart.Start = chart.End.Add(-c.Duration("d"))
	chart.SampleInterval = c.Duration("d") / 600

	query := "sum(rate(traefik_backend_requests_total[5m])) by (backend)"

	data, err := util.GetChartData(query, chart.Start, chart.End)

	if err != nil {
		log.Fatalln("Error getting chart data:", err)
	}

	chart.Series = data

	p.DrawTimeSeries(chart)

}
