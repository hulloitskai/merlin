package models

// A Filing represents a company filing on EDGAR.
type Filing struct {
	Type   string `json:"type"`
	Desc   string `json:"description"`
	Date   string `json:"date"`
	AccNum string `json:"accessionNumber"`
}

// Filings are a set of Filing objects.
type Filings []*Filing

// FilingResults are the results of a filings search query.
type FilingResults struct {
	CIK     string `json:"CIK"`
	Filings `json:"filings"`
}
