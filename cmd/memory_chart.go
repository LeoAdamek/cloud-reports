package cmd

import (
	"encoding/json"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"net/url"
	"reports/printer"
	"strconv"
	"time"
)

func memoryChart(c *cli.Context) {

	data := getMemoryChartData()

	p := getPrinter(c)

	chart := printer.TimeSeriesChart{
		End: time.Now().Truncate(time.Hour),
		SampleInterval: 144*time.Second,
		Series: data,
		Title: "",
		Unit: "GiB",
	}

	chart.Start = chart.End.AddDate(0,0,-1)

	p.DrawTimeSeries(chart)

}

func getMemoryChartData() [][]float64 {

	end := time.Now().Truncate(time.Hour)
	start := end.AddDate(0,0,-1)
	queryStr := "sum(node_memory_MemAvailable_bytes) by (instance) / 6"


	urlStr := "https://prometheus.adamek.me/api/v1/query_range"

	u, err := url.Parse(urlStr)

	if err != nil {
		log.Fatalln("Invalid Promethus Host:", err)
	}

	q := u.Query()
	q.Set("query", queryStr)
	q.Set("start", strconv.FormatInt(start.Unix(), 10))
	q.Set("end", strconv.FormatInt(end.Unix(),10))
	q.Set("step", "144")

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

			series[sidx][i] /= 1<<30
		}
	}


	return series

}