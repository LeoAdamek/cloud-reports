package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const hostStr = "https://prometheus.adamek.me/api/v1/query_range"

type QueryResponse struct {
	Data *ResponseData`json:"data"`
}

type ResponseData struct {
	Result []Result `json:"result"`
}

type Result struct {
	Values []Value `json:"values"`
}

type Value struct {
	T time.Time
	V float64
}

func QueryRange(query string, start, end time.Time, interval time.Duration) (*QueryResponse, error) {


	u, _ := url.Parse(hostStr)

	q := u.Query()

	q.Set("query", query)
	q.Set("start", strconv.FormatInt(start.Unix(), 10))
	q.Set("end"  , strconv.FormatInt(end.Unix(), 10))
	q.Set("step",  strconv.FormatInt(int64(interval.Truncate(time.Second).Seconds()), 10))

	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())

	if err != nil {
		return nil, err
	}

	data := &QueryResponse{}

	err = json.NewDecoder(res.Body).Decode(data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetChartData(query string, start, end time.Time) ([][]float64, error) {

	step := end.Sub(start) / 600

	data, err := QueryRange(query, start, end, step)

	if err != nil {
		return [][]float64{{}}, err
	}

	series := make([][]float64, len(data.Data.Result))

	for j, s := range data.Data.Result {
		sl := make([]float64, len(s.Values))

		for i, v := range s.Values {
			sl[i] = v.V
		}

		series[j] = sl
	}

	return series, nil
}


func (v *Value) UnmarshalJSON(in []byte) error {
	var content []interface{}

	err := json.Unmarshal(in, &content)

	if err != nil {
		return err
	}

	ts, ok := content[0].(float64)

	if !ok {
		return errors.New(fmt.Sprintf("first element of a sample must be float64, got %#+v", content[0]))
	}

	v.T = time.Unix(int64(math.Floor(ts)), 0)

	iv, ok := content[1].(string)

	if !ok {
		return errors.New("second element must be a string")
	}

	val, err := strconv.ParseFloat(iv, 64)

	if err != nil {
		return err
	}

	v.V = val

	return nil
}