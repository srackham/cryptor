# Cryptor

<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/srackham/cryptor)](https://goreportcard.com/report/github.com/srackham/cryptor) -->

Cryptor reports crypto currency portfolio statistics (current and historical).

- Cryptor is a CLI application written in Go.
- Cryptor handles multiple asset portfolios and historic positions.
- Cryptor tracks the value and performance of crypto assets, it does not communicate or integrate with blockchains or wallets.

## Quick Start
To compile and install the `cryptor` executable you will first need to [Download and install the Go Programming Language](https://go.dev/doc/install).

Next run this command:

    go install TODO

Test it by running:

    cryptor help


## Implementation Notes
- To ensure minimal use of Web currency APIs (which are often throttled) Portfolio evaluations, crypto prices and exchange rates are cached to the local `$HOME/.cryptor` directory. Normally, unless the `-refresh` option is used, current asset prices and exchange rates will only be fetched over the Internet once per day.

- Asset values are saved in USD.

- Current portfolio valuations are saved to the `$HOME/.cryptor/history.json` file. Historic valuations (using the `-date` option) are not saved.

- When a portfolio valuation is saved to the history file it is date stamped with the current date. If the `-refresh` option is used the portfolio is also time stamped (otherwise the time is left blank).

- Dates are saved as YYYY-DD-MM formatted strings; times are saved as HH:MM:SS formatted strings. Local dates and times are used.

- Use the `-currency` option to display values in non-USD currencies.

- All non-USD valuations are based on the current exchange rate against the USD irrespective of the `-date` option.
