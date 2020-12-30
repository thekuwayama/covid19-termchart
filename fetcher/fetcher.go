package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

type Fetcher struct {
	httpClient HttpClient
}

type HttpClient interface {
	Get(string) (*http.Response, error)
}

func New(c HttpClient) *Fetcher {
	f := &Fetcher{httpClient: c}
	return f
}

func (fetcher *Fetcher) Do(s string) (string, error) {
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
