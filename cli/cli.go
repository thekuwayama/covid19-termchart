package cli

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/gat/render"
	"github.com/thekuwayama/covid19-termchart/fetcher"
	"github.com/thekuwayama/covid19-termchart/plotter"
)

const OpenDataUrl = "https://www3.nhk.or.jp/n-data/opendata/coronavirus/nhk_news_covid19_domestic_daily_data.csv"

func Run() {
	f := fetcher.New(http.DefaultClient)
	csv, err := f.Do(OpenDataUrl)
	if err != nil {
		log.Fatal(err)
	}

	b, err := plotter.Plot(csv)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	err = printImage(img)
	if err != nil {
		log.Fatal(err)
	}
}

func printImage(img image.Image) error {
	iterm := &render.ITerm{}
	return iterm.Render(os.Stdout, img)
}
