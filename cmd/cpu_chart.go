package cmd

import (
	"encoding/json"
	"github.com/urfave/cli"
	"log"
	"math"
	"net/http"
	"net/url"
	"reports/printer"
	"strconv"
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


	data := getCpuChartData(chart.Start, chart.End, chart.SampleInterval, c.StringSlice("group"))

	chart.Series = data

	p.DrawTimeSeries(chart)
	
}

func getCpuChartData(start, end time.Time, interval time.Duration, groupBy []string) [][]float64 {

	if len(groupBy) == 0 {
		groupBy = []string{"instance"}
	}

	queryStr := "sum(rate(node_cpu_seconds_total{mode!=\"idle\"}[30m])) by ("+strings.Join(groupBy, ",")+")"
	

	urlStr := "https://prometheus.adamek.me/api/v1/query_range"
	
	u, err := url.Parse(urlStr)
	
	if err != nil {
		log.Fatalln("Invalid Promethus Host:", err)
	}

	q := u.Query()
	q.Set("query", queryStr)
	q.Set("start", strconv.FormatInt(start.Unix(), 10))
	q.Set("end", strconv.FormatInt(end.Unix(),10))
	q.Set("step", strconv.FormatInt(int64(math.Floor(interval.Seconds())), 10))
	
	u.RawQuery = q.Encode()
	
	res, err := http.Get(u.String())
	
	if err != nil {
		log.Println("Request Error:", err)
	}
	
	data := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&data)
	
	if err != nil {
		log.Println("Decode Error:", err)
	}
	
	results := data["data"].(map[string]interface{})["result"].([]interface{})
	
	series := make([][]float64, len(results))
	
	for sidx, ser := range results {
		qs := ser.(map[string]interface{})
		sv := qs["values"].([]interface{})
		
		series[sidx] = make([]float64, len(sv))
		
		log.Printf("Got series length: %d", len(sv))
		
		for i, vs := range sv {
			vi := vs.([]interface{})
			series[sidx][i], _ = strconv.ParseFloat(vi[1].(string), 64)

			series[sidx][i] *= 100
		}
	}
	
	
	return series
	
}