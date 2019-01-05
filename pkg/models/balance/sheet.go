package balance

import (
	"fmt"
)

// A Sheet is a filing that contains information about company balances.
type Sheet struct {
	CIK      string
	AccNum   string     `json:"accessionNumber"`
	Date     string     `json:"date"`
	Sections SectionMap `json:"sections"`
}

// NewSheet returns a new Sheet.
func NewSheet(cik, accNum, date string) *Sheet {
	return &Sheet{
		CIK:      cik,
		AccNum:   accNum,
		Date:     date,
		Sections: make(SectionMap),
	}
}

// AddItem adds an Item to the Sheet, in the section corresponding to sec.
func (s *Sheet) AddItem(sec Section, i *Item) {
	s.Sections[sec] = append(s.Sections[sec], i)
}

func (s *Sheet) String() string {
	return fmt.Sprintf("Sheet{CIK: %s, AccNum: %s, Date: %s, Sections: %v}",
		s.CIK, s.AccNum, s.Date, s.Sections)
}

// Sheets are a set of Sheet objects.
type Sheets []*Sheet
