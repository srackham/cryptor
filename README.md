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


## Commands
TODO


## Valuate Command
The _valuate_ command calculates portfolio asset values in USD (or some other specified currency).

### Syntax

    cryptor valuate [OPTION]...

### Options
[Common command options](#common-command-options) plus:

    -aggregate
    -currency CURRENCY
    -date DATE
    -refresh

- Use the `-currency` option to display values in non-USD currencies.
- All non-USD valuations are based on the current exchange rate against the USD.
- The `-date DATE` option specifies the date of the crypto prices; the default date is today's date.

### Examples
TODO

## Implementation Notes
- To ensure minimal use of Web currency APIs (which are often throttled) Portfolio valuations, crypto prices and exchange rates are cached to the local `$HOME/.cryptor` directory. Use the `-refresh` option to bypass the cache and fetch the latest asset prices and exchange rates.

- Current portfolio valuations are saved to the `$HOME/.cryptor/history.json` file. Historic valuations (using the `-date` option) are not saved.

- Asset values are saved in USD.

- Saved portfolio valuations are date stamped. Current valuations fetched with the `-refresh` option are also time stamped (otherwise the time is left blank).

- Dates are specified and saved as YYYY-DD-MM formatted strings; times are saved as HH:MM:SS formatted strings. Local dates and times are used.


## Tips
-