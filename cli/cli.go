package cli

import (
	"log"
	"net/http"

	"github.com/thekuwayama/covid19-termchart/fetcher"
	"github.com/thekuwayama/covid19-termchart/termui"
)

const OpenDataUrl = "https://www3.nhk.or.jp/n-data/opendata/coronavirus/nhk_news_covid19_domestic_daily_data.csv"

func Run() {
	f := fetcher.New(http.DefaultClient)
	csv, err := f.Do(OpenDataUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = termui.Plot(csv)
	if err != nil {
		log.Fatal(err)
	}
}
