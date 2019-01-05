# merlin

[![Github: Tag][tag-img]][tag]
[![Travis: Build][travis-img]][travis]
[![Go Report Card][grp-img]][grp]

_A system for accessing company finance data from
[EDGAR](https://www.sec.gov/edgar/aboutedgar.htm)._

## API

| Endpoint               | Description                                              |
| ---------------------- | -------------------------------------------------------- |
| `/`                    | API server information.                                  |
| `/sheets/:cik/:accNum` | Balance sheet data for a given CIK and accession number. |

> When accessing `merlin` on the production server at
> https://merlin.stevenxie.me, all API requests must be prefixed with `/api`.
>
> For example, to access the API root (`/`), one would visit `/api/` instead.

## TODOs

- [x] Add an endpoint for parsing a balance sheet corresponding to a CIK +
      accNum.
- [ ] Add an endpoint for mapping tickers to CIKs.
- [ ] Add an endpoint for listing accNums for a given CIK.
- [ ] Add result caching using Redis to improve speeds.

[tag]: https://github.com/stevenxie/merlin/releases
[tag-img]: https://img.shields.io/github/tag/stevenxie/merlin.svg
[travis]: https://travis-ci.com/stevenxie/merlin
[travis-img]: https://travis-ci.com/stevenxie/merlin.svg?branch=master
[grp]: https://goreportcard.com/report/github.com/stevenxie/merlin
[grp-img]: https://goreportcard.com/badge/github.com/stevenxie/merlin
