package scrape

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
	"github.com/stevenxie/merlin/pkg/models/balance"
	ess "github.com/unixpickle/essentials"
	"golang.org/x/text/runes"
)

// ScrapeBalanceSheets scrapes EDGAR for the balance sheets for the filing
// associated with cik and accNum.
func (s *Scraper) ScrapeBalanceSheets(cik, accNum string) (
	balance.Sheets, error) {
	doc, err := s.getBSDoc(cik, accNum)
	if err != nil {
		return nil, ess.AddCtx("scrape: getting balance sheet document", err)
	}

	sel := doc.Find(".report").Children()
	if sel.Length() != 1 {
		return nil, errors.New("scrape: cannot find report table")
	}

	// Parse date row.
	var (
		rows    = sel.Children()
		dateRow = rows.First()
		dates   []string
	)
	dateRow.Children().EachWithBreak(func(_ int, cell *gq.Selection) bool {
		if _, ok := cell.Attr("colspan"); ok { // skip cell if it has 'colspan'
			return true
		}

		date := cell.Text()
		if date == "" {
			err = errors.New("date cell has no text")
			return false
		}

		dates = append(dates, date)
		return true
	})
	if err != nil {
		return nil, ess.AddCtx("scrape: parsing date row", err)
	}

	// Early return on empty date row.
	if len(dates) == 0 {
		return nil, nil
	}

	// Initialize balance sheets.
	sheets := make(balance.Sheets, len(dates))
	for i, date := range dates {
		sheets[i] = balance.NewSheet(cik, accNum, date)
	}

	// Trim title row.
	rows.Nodes = rows.Nodes[1:]

	const badSection = -1
	var sec balance.Section = badSection
	rows.EachWithBreak(func(i int, row *gq.Selection) bool {
		// Check if the row is a section header (marked by the presence of a
		// <strong> element).
		if row.Find("strong").Length() > 0 {
			if sec, err = parseBSHeaderRow(row, sheets); err != nil {
				goto catch
			}
			return true
		}

		// Skip row if section is invalid.
		if sec == badSection {
			return true
		}
		// Skip row if it is an extraneous header (selector: '.rh').
		if class, _ := row.Attr("class"); class == "rh" {
			return true
		}

		if err = parseBSRow(row, sheets, sec); err != nil {
			goto catch
		}
		return true

	catch:
		ess.AddCtxTo(fmt.Sprintf("row %d", i), &err)
		return false
	})

	return sheets, ess.AddCtx("scrape: parsing balance sheet rows", err)
}

// parseBSHeaderRow parses a section header row into a balance.Section.
//
// Returns -1 if the parsed header does not match any known balance.Section.
func parseBSHeaderRow(row *gq.Selection, sheets balance.Sheets) (
	balance.Section, error,
) {
	text := strings.TrimSpace(row.Children().First().Text())
	if text == "" {
		return 0, errors.New("section name cell is empty")
	}
	text = strings.ToLower(text)

	// Remove characters from cutset.
	cutset := "'’"
	set := runes.Predicate(func(r rune) bool {
		return strings.ContainsRune(cutset, r)
	})
	transformer := runes.Remove(set)
	text = transformer.String(text)

	const notFound = -1
	if strings.Contains(text, "assets") {
		if strings.Contains(text, "non-current") {
			return balance.SecNonCurrAssets, nil
		}
		if strings.Contains(text, "current") {
			return balance.SecCurrAssets, nil
		}
		return notFound, nil
	}
	if strings.Contains(text, "liabilities") {
		if strings.Contains(text, "non-current") {
			return balance.SecNonCurrLiabilities, nil
		}
		if strings.Contains(text, "current") {
			return balance.SecCurrLiabilities, nil

		}
		return notFound, nil
	}
	if strings.Contains(text, "shareholders equity") ||
		strings.Contains(text, "stockholders equity") {
		return balance.SecStockholdersEquity, nil
	}
	return notFound, nil
}

// parseBSRow parses row into balance.Items, and stores it in sheets.
func parseBSRow(row *gq.Selection, sheets balance.Sheets,
	sec balance.Section) error {
	var (
		err      error
		template balance.Item
	)
	row.Children().EachWithBreak(func(i int, cell *gq.Selection) bool {
		// Is name cell.
		if i == 0 {
			if template.Name = strings.TrimSpace(cell.Text()); template.Name == "" {
				err = errors.New("name cell is empty")
				return false
			}
			// Normalize quotations.
			template.Name = strings.Replace(template.Name, "’", "'", -1)

			// Parse trailing note tag.
			// TODO: Figure out how to parse multiple note tags.
			const marker = "(note "
			var (
				name   = strings.ToLower(template.Name)
				tagi   = strings.LastIndex(name, marker)
				closei = strings.LastIndexByte(name, ')')
			)
			if (tagi == -1) || (closei == -1) {
				return true // skip
			}

			var (
				sid = name[tagi+len(marker) : closei]
				id  int
			)
			if id, err = strconv.Atoi(sid); err != nil {
				ess.AddCtxTo("parsing note ID as integer", &err)
				return false
			}
			template.Notes = append(template.Notes, id)
			template.Name = strings.TrimSpace(template.Name[:tagi])
			return true
		}

		// Is value cell.
		item := template // copy template
		item.Value = strings.TrimSpace(cell.Text())
		sheets[i-1].AddItem(sec, &item)
		return true
	})

	return err
}

// getBSDoc returns a gq.Document representing to the balance sheet of the
// filing that corresponds to cik and accNum.
func (s *Scraper) getBSDoc(cik, accNum string) (*gq.Document, error) {
	uri, err := s.deriveBSURL(cik, accNum)
	if err != nil {
		return nil, ess.AddCtx("deriving balance sheet URL", err)
	}
	return s.readDocumentAt(uri)
}

// deriveBSURL looks up the filing associated with cik and accNum on the
// EDGAR viewer, and attempts to derive the corresponding balance sheet URL.
func (s *Scraper) deriveBSURL(cik, accNum string) (uri string,
	err error) {
	// Perform request.
	u := makeViewerURL(cik, accNum)
	res, err := s.Client.Get(u)
	if err != nil {
		return "", ess.AddCtx("getting EDGAR viewer", err)
	}
	defer res.Body.Close()

	// Parse response body for sidebar.
	doc, err := gq.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", ess.AddCtx("parsing response body with goquery", err)
	}
	if err = res.Body.Close(); err != nil {
		return "", ess.AddCtx("closing response body", err)
	}

	sel := doc.Find("#menu_cat2")
	if sel.Length() != 1 {
		return "", errors.New("could not not find financial statements sidebar " +
			"item")
	}

	sel = sel.Parent().Find("ul")
	if sel.Length() != 1 {
		return "", errors.New("could not find list of financial statements")
	}

	// Check to see if any of the sidebar items match the query terms.
	var id string
	sel.Children().EachWithBreak(func(i int, item *gq.Selection) bool {
		text := item.Text()
		if text == "" {
			return true // skip
		}

		// Process text.
		text = strings.ToLower(text)
		if !strings.Contains(text, "balance sheet") ||
			strings.Contains(text, "parenthetical") {
			return true // skip
		}

		var ok bool
		if id, ok = item.Attr("id"); !ok {
			err = fmt.Errorf("item %d: no id", i)
		}
		return false
	})
	if err != nil {
		return "", ess.AddCtx("parsing financial statements", err)
	}

	if id == "" {
		return "", errors.New("could not find balance sheets item in sidebar")
	}
	id = strings.ToUpper(id)
	flatAccNum := strings.Replace(accNum, "-", "", -1)
	return fmt.Sprintf("%s/%s/%s/%s.htm", edgarDataURL, cik, flatAccNum, id), nil
}
