package scrape

import "net/url"

const (
	edgarDataURL   = "https://www.sec.gov/Archives/edgar/data"
	edgarViewerURL = "https://www.sec.gov/cgi-bin/viewer"
)

func makeViewerURL(cik, accNum string) string {
	u, err := url.Parse(edgarViewerURL)
	if err != nil {
		panic(err)
	}

	params := u.Query()
	params.Set("action", "view")
	params.Set("cik", cik)
	params.Set("accession_number", accNum)
	u.RawQuery = params.Encode()

	return u.String()
}
