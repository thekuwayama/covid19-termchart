package termui

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/gat/render"
	chart "github.com/wcharczuk/go-chart/v2"
	"golang.org/x/xerrors"
)

type chartValue struct {
	x []float64
	y []float64
}

func parse(csv string) (*chartValue, error) {
	x := []float64{}
	y := []float64{}
	lines := strings.Split(strings.ReplaceAll(csv, "\r\n", "\n"), "\n")[1:] // skip header
	for _, l := range lines {
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

func Plot(csv string) error {
	xy, err := parse(csv)
	if err != nil {
		return xerrors.Errorf("Failed to parse csv: %+w", err)
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := chart.TimeFromFloat64(typed)
				return fmt.Sprintf("%d/%d/%d", typedDate.Year(), typedDate.Month(), typedDate.Day())
			},
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

	f, err := ioutil.TempFile("", "output*.png")
	if err != nil {
		return xerrors.Errorf("Failed to touch tmp file: %+w", err)
	}
	defer f.Close()

	graph.Render(chart.PNG, f)
	f.Seek(0, 0)
	img, err := readImage(f)
	if err != nil {
		return xerrors.Errorf("Failed to read image file: %+w", err)
	}

	err = printImage(img)
	if err != nil {
		return xerrors.Errorf("Failed to print image: %+w", err)
	}
	return nil
}

func readImage(f *os.File) (image.Image, error) {
	var r io.Reader
	r = f
	return png.Decode(r)
}

func printImage(img image.Image) error {
	iterm := &render.ITerm{}
	return iterm.Render(os.Stdout, img)
}
