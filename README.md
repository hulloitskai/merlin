# merlin

[![Github: Tag][tag-img]][tag]
[![Travis: Build][travis-img]][travis]
[![Go Report Card][grp-img]][grp]

_A system for accessing company finance data from
[EDGAR](https://www.sec.gov/edgar/aboutedgar.htm)._

| Deployment                  | API                             |
| --------------------------- | ------------------------------- |
| https://merlin.stevenxie.me | https://merlin.stevenxie.me/api |

## API

| Endpoint                      | Description                                                     |
| ----------------------------- | --------------------------------------------------------------- |
| `/`                           | API server information.                                         |
| `/filings/:ticker/`           | Company filings (CIK and accession numbers) for a given ticker. |
| `/filings/:ticker/latest/10k` | Latest 10-K filing for a given ticker.                          |
| `/sheets/:cik/:acc-num`       | Balance sheet data for a given CIK and accession number.        |
| `/notes/:cik/:acc-num`        | Financial notes for a given CIK and accession number.           |

> When accessing `merlin` on the production server at
> https://merlin.stevenxie.me, all API requests must be prefixed with `/api`.
>
> For example, to access the API root (`/`), one would visit `/api/` instead.

### Examples

#### Company Filings:

`GET` https://merlin.stevenxie.me/api/filings/MSFT/

```jsonc
{
  "CIK": "0000789019",
  "filings": [
    {
      "type": "8-K",
      "description": "Current report",
      "date": "2018-11-29",
      "accessionNumber": "0001193125-18-337951"
    }
    // ...
  ]
}
```

`GET` https://merlin.stevenxie.me/api/filings/MSFT/latest/10k

```jsonc
{
  "CIK": "0000789019",
  "filing": {
    "type": "10-K",
    "description": "Annual report [Section 13 and 15(d), not S-K Item 405]",
    "date": "2016-07-28",
    "accessionNumber": "0001193125-16-662209"
  }
}
```

#### Balance Sheets:

`GET` https://merlin.stevenxie.me/api/sheets/1318605/0001564590-18-002956

```jsonc
[
  {
    "CIK": "1318605",
    "accessionNumber": "0001564590-18-002956",
    "date": "Dec. 31, 2017",
    "sections": {
      "currentAssets": [
        {
          "name": "Cash and cash equivalents",
          "value": "$ 3,367,914"
        }
        // ...
      ]
    }
  }
  // ...
]
```

#### Financial Notes:

`GET` https://merlin.stevenxie.me/api/notes/1318605/0001564590-18-002956

```jsonc
[
  {
    "id": 1,
    "title": "Overview",
    "link": "https://www.sec.gov/Archives/edgar/data/1318605/000156459018002956/R9.htm"
  }
  // ...
]
```

<br />

## TODOs

- [x] Add an endpoint for parsing a balance sheet corresponding to a CIK +
      accession number.
- [x] Add an endpoint for parsing financial notes corresponding to a CIK +
      accession number.
- [x] Add sample frontend.
- [x] Add an endpoint for mapping tickers to CIKs.
- [x] Add an endpoint for listing accNums for a given CIK.
- [ ] Try a new information extraction strategy using XLSX financial reports
      ([like this](https://www.sec.gov/Archives/edgar/data/789019/000119312516662209/)).
- [ ] Add result caching using Redis to improve speeds.

## Further Notices

- Not all balance sheets can be parsed properly using the current strategy,
  see https://www.sec.gov/Archives/edgar/data/789019/000119312516662209/R5.htm
  (notice the inconsistent column count). This can be caught with either
  using right-delta indexes when counting date columns, or with XLSX parsing
  to gain more contextual information about the table.

[tag]: https://github.com/stevenxie/merlin/releases
[tag-img]: https://img.shields.io/github/tag/stevenxie/merlin.svg
[travis]: https://travis-ci.com/stevenxie/merlin
[travis-img]: https://travis-ci.com/stevenxie/merlin.svg?branch=master
[grp]: https://goreportcard.com/report/github.com/stevenxie/merlin
[grp-img]: https://goreportcard.com/badge/github.com/stevenxie/merlin
