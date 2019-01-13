package scrape

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/stevenxie/merlin/api/pkg/models"
	ess "github.com/unixpickle/essentials"
)

// ScrapeFilings scrapes EDGAR for company filings associated with the given
// ticker.
func (s *Scraper) ScrapeFilings(ticker string) (*models.FilingResults, error) {
	u, err := url.Parse(edgarBrowseURL)
	if err != nil {
		panic(err)
	}

	params := u.Query()
	params.Set("action", "getcompany")
	params.Set("owner", "exclude")
	params.Set("count", "100")
	params.Set("output", "xml")
	params.Set("CIK", ticker)
	u.RawQuery = params.Encode()

	res, err := s.Client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	ctype := res.Header.Get("Content-Type")
	if ctype != "application/xml" {
		return nil, fmt.Errorf("scrape: server did not respond with XML "+
			"(Content-Type='%s')", ctype)
	}

	// Decode XML from response body.
	var result struct {
		CIK     string `xml:"companyInfo>CIK"`
		Filings []struct {
			DateFiled  string `xml:"dateFiled"`
			FilingHREF string `xml:"filingHREF"`
			FormName   string `xml:"formName"`
			Type       string `xml:"type"`
		} `xml:"results>filing"`
	}
	dec := xml.NewDecoder(res.Body)
	dec.CharsetReader = charset.NewReaderLabel
	if err = dec.Decode(&result); err != nil {
		return nil, ess.AddCtx("scrape: decoding response body as XML", err)
	}

	fr := &models.FilingResults{CIK: result.CIK}
	for i, filing := range result.Filings {
		slashi := strings.LastIndexByte(filing.FilingHREF, '/')
		if slashi == -1 {
			return nil, fmt.Errorf("scrape: filing %d: could not find path slashes "+
				"in filing HREF", i)
		}
		accNum := filing.FilingHREF[slashi+1:]
		accNum = strings.TrimSuffix(accNum, "-index.htm")

		fr.Filings = append(fr.Filings, &models.Filing{
			AccNum: accNum,
			Date:   filing.DateFiled,
			Desc:   filing.FormName,
			Type:   filing.Type,
		})
	}

	err = res.Body.Close()
	return fr, ess.AddCtx("scrape: closing response body", err)
}
