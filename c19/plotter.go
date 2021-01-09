package c19

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
	"golang.org/x/xerrors"
)

type chartValue struct {
	x float64
	y float64
}

func parse(csv string, days int) ([]chartValue, error) {
	lines := strings.Split(strings.ReplaceAll(csv, "\r\n", "\n"), "\n")[1:] // skip header
	res := make([]chartValue, 0, len(lines))
	for i, l := range lines {
		if i > days {
			break
		}

		xy := strings.Split(l, ",")
		if len(xy) != 5 {
			continue
		}

		t, err := time.Parse("2006/1/2", xy[0])
		if err != nil {
			return nil, err
		}

		x := float64(t.UnixNano())
		y, err := strconv.ParseFloat(xy[1], 64)
		if err != nil {
			return nil, err
		}

		res = append(res, chartValue{x: x, y: y})
	}

	return res, nil
}

func calcWeeklyAverage(cvs []chartValue) []chartValue {
	l := len(cvs)
	res := make([]chartValue, l)
	for i := 0; i < l; i++ {
		res[i].x = cvs[i].x
		if i < 6 {
			res[i].y = 0
			continue
		}

		res[i].y = (cvs[i-6].y + cvs[i-5].y + cvs[i-4].y + cvs[i-3].y + cvs[i-2].y + cvs[i-1].y + cvs[i].y) / 7
	}

	return res
}

func Plot(csv string, days int) ([]byte, error) {
	xy, err := parse(csv, days)
	if err != nil {
		return nil, xerrors.Errorf("Failed to parse csv: %+w", err)
	}

	weekly := calcWeeklyAverage(xy)
	l := len(xy)
	if l != len(weekly) {
		return nil, xerrors.Errorf("Failed to calculate weekly average")
	}

	x := make([]float64, l)
	y := make([]float64, l)
	w := make([]float64, l)
	for i := 0; i < l; i++ {
		x[i] = xy[i].x
		y[i] = xy[i].y
		w[i] = weekly[i].y
	}

	graph := chart.Chart{
		Title: "Number of infected people in Japan (NHK summary)",
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := chart.TimeFromFloat64(typed)
				return fmt.Sprintf("%d/%d/%d", typedDate.Year(), typedDate.Month(), typedDate.Day())
			},
		},
		YAxis: chart.YAxis{
			Name:           "Number of infected people in Japan /day & weekly avg",
			ValueFormatter: chart.IntValueFormatter,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: x,
				YValues: y,
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(255),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(0),
				},
				XValues: x,
				YValues: w,
			},
		},
	}

	bb := bytes.NewBufferString("")
	err = graph.Render(chart.PNG, bb)
	if err != nil {
		return nil, xerrors.Errorf("Failed to render image: %+w", err)
	}

	return bb.Bytes(), nil
}
