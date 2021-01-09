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
	x []float64
	y []float64
}

func parse(csv string, days int) (*chartValue, error) {
	x := []float64{}
	y := []float64{}
	lines := strings.Split(strings.ReplaceAll(csv, "\r\n", "\n"), "\n")[1:] // skip header
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
		x = append(x, float64(t.UnixNano()))

		f, err := strconv.ParseFloat(xy[1], 64)
		if err != nil {
			return nil, err
		}
		y = append(y, f)
	}

	return &chartValue{x: x, y: y}, nil
}

func Plot(csv string, days int) ([]byte, error) {
	xy, err := parse(csv, days)
	if err != nil {
		return nil, xerrors.Errorf("Failed to parse csv: %+w", err)
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
			Name:           "Number of infected people in Japan /day",
			ValueFormatter: chart.IntValueFormatter,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: xy.x,
				YValues: xy.y,
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
