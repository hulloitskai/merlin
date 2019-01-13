package balance

import (
	"fmt"
)

// A Section corresponds to a section on a Sheet.
type Section int8

// Constants corresponding to Sheet sections.
const (
	SecCurrAssets Section = iota
	SecNonCurrAssets
	SecCurrLiabilities
	SecNonCurrLiabilities
	SecStockholdersEquity
)

func (sec Section) String() string {
	switch sec {
	case SecCurrAssets:
		return "Current Assets"
	case SecNonCurrAssets:
		return "Non-current Assets"
	case SecCurrLiabilities:
		return "Current Liabilities"
	case SecNonCurrLiabilities:
		return "Non-current Liabilities"
	case SecStockholdersEquity:
		return "Stockholder's Equity"
	default:
		err := fmt.Errorf("balance: no such Section '%d'", sec)
		panic(err)
	}
}

// MarshalText implements format.TextMarshaler for sec.
func (sec Section) MarshalText() ([]byte, error) {
	var val string
	switch sec {
	case SecCurrAssets:
		val = "currentAssets"
	case SecNonCurrAssets:
		val = "nonCurrentAssets"
	case SecCurrLiabilities:
		val = "currentLiabilities"
	case SecNonCurrLiabilities:
		val = "nonCurrentLiabilities"
	case SecStockholdersEquity:
		val = "stockholdersEquity"
	default:
		err := fmt.Errorf("balance: no such Section '%d'", sec)
		panic(err)
	}
	return []byte(val), nil
}

// SectionMap is a map of Sections to Items.
type SectionMap map[Section]Items
