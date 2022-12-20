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


## Glossary
- _Allocation_: TODO
- _Amount_: TODO
- _Asset_: TODO
- _Cost_: TODO
- _Currency_: TODO
- _Portfolio_: TODO
- _Price_: TODO
- _Value_: TODO


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
- Valuations with non-USD currencies are based on current exchange rates against the USD.
- The `-date DATE` option specifies the date of the crypto prices; the default date is today's date.

### Examples
TODO

## Implementation Notes
- Crypto prices and exchange rates are cached locally to the cryptor configuration directory (defaults to `$HOME/.cryptor`). Unless the `-force` option is specified, updates are only fetched over the Internet when they are not found in the local cache files. Caching ensure minimal use of Web APIs which require an Internet connection, can be slow and are sometimes throttled.

- By convention, crypto and currency symbols are converted to upper case.

- The configuration directory can be specified with the `-confdir` option.

- Portfolio valuations are saved to the `$HOME/.cryptor/valuations.json` file:
    * Valuations of past dates (using the `-date` option) are not saved unless the `-force` option is also used.
    * Valuations for the current day always update the valuations history.
    * A current valuation made with the `-force` option ensures that the latest crypto prices are used.

- All values are saved in USD.
- Saved portfolio valuations are date stamped.
- Dates are specified and saved as YYYY-DD-MM formatted strings in the local time zone.


## Tips
-