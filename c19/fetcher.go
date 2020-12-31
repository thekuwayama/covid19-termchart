package c19

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

const OpenDataUrl = "https://www3.nhk.or.jp/n-data/opendata/coronavirus/nhk_news_covid19_domestic_daily_data.csv"

type Fetcher struct {
	httpClient HttpClient
}

type HttpClient interface {
	Get(string) (*http.Response, error)
}

func NewFetcher(c HttpClient) *Fetcher {
	f := &Fetcher{httpClient: c}
	return f
}

func (fetcher *Fetcher) Fetch(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", xerrors.Errorf("Failed to parse URL: %+w", err)
	}

	res, err := fetcher.httpClient.Get(u.String())
	if err != nil {
		return "", xerrors.Errorf("Failed to do the HTTP GET: %+w", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", xerrors.Errorf("Failed to read HTTP body: %+w", err)
	}
	return string(b), nil
}
