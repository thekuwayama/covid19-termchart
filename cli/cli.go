package cli

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/gat/render"
	"github.com/thekuwayama/covid19-termchart/c19"
)

var days = flag.Int("day", 365, "period to aggregate")

func Run() {
	flag.Parse()

	f := c19.NewFetcher(http.DefaultClient)
	csv, err := f.Fetch(c19.OpenDataUrl)
	if err != nil {
		log.Fatal(err)
	}

	b, err := c19.Plot(csv, *days)
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
