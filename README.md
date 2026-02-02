# go-evatr

[![Go Reference](https://pkg.go.dev/badge/github.com/hostwithquantum/go-evatr.svg)](https://pkg.go.dev/github.com/hostwithquantum/go-evatr)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/hostwithquantum/go-evatr/badge)](https://scorecard.dev/viewer/?uri=github.com/hostwithquantum/go-evatr)

Go client library for the eVatR API (German VAT ID validation system) of the BZSt.

## Features

- Simple and qualified VAT ID validation
- EU member state information and VIES availability
- Type-safe error handling with status codes
- Context-aware API calls

## Installation

```bash
go get github.com/hostwithquantum/go-evatr
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/hostwithquantum/go-evatr) - Full API documentation
- [Examples](./examples/) - Working code examples
- [Error codes](errors.go) - All status codes and their meanings

### Optional Configuration

```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: customTransport,
}

client := evatr.NewClient(
    evatr.WithBaseURL("https://custom.api.url"),
    evatr.WithTimeout(60 * time.Second),
    evatr.WithHTTPClient(httpClient),
)
```

### Usage

The API _advertises_ a daily maintenance window from 23:00 - 5:00 (local). Run potential jobs during the workday to avoid issues — see our dependabot and workflow configuration for examples.

## License

[mpl-2.0](./LICENSE)

## Links

- [eVatR API Documentation](https://api.evatr.vies.bzst.de/api-docs)
- [BZSt (Federal Central Tax Office)](https://www.bzst.de/DE/Unternehmen/Identifikationsnummern/Umsatzsteuer-Identifikationsnummer/AuslaendischeUSt-IdNr/auslaendische_ust_idnr_node.html)
- [VIES VAT Validation](https://ec.europa.eu/taxation_customs/vies/)
