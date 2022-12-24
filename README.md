# Cryptor

Cryptor valuates crypto currency asset portfolios.

- Cryptor can process multiple asset portfolios and historic valuations.
- Cryptor tracks the values and performance of crypto assets
- Cryptor uses publicly available crypto prices and exchange rates, it does not communicate or integrate with blockchains or wallets.
- Cryptor is a CLI application written in Go.

## Quick Start
Install `cryptor` with this command (prerequisite:
[the Go Programming Language](https://go.dev/doc/install)):

    go install github.com/srackham/cryptor@latest

Install an example portfolios configuration file using the `cryptor init` command. For example:

```
$ cryptor init
creating configuration directory: "/home/srackham/.cryptor"
installing example portfolios file: "/home/srackham/.cryptor/portfolios.yaml"
```

Edit the YAML portfolios configuration file (`$HOME/.cryptor/portfolios.yaml`) to match your own:

```yaml
# Example cryptor portfolio configuration file

- name:  personal
  notes: Personal portfolio
  cost: $10,000.00 NZD
  assets:
    BTC: 0.5
    ETH: 2.5
    USDC: 100

- name:  joint
  assets:
      BTC: 0.5
      ETH: 2.5

# Minimal portfolio
- assets:
      BTC: 0.25
```

Use the `cryptor valuate` command to value the portfolios. For example:

```
$ cryptor valuate

NAME:  personal
NOTES: Personal portfolio
DATE:  2022-12-22
VALUE: 11574.20 USD
COST:  6319.93 USD
GAINS: 5254.27 (83.14%)
XRATE:
            AMOUNT            VALUE   PERCENT            PRICE
BTC         0.5000      8430.65 USD    72.84%     16861.30 USD
ETH         2.5000      3043.78 USD    26.30%      1217.51 USD
USDC      100.0000        99.77 USD     0.86%         1.00 USD
```

Run the `cryptor help` command to view all the commands and command options:

```
$ cryptor help

Cryptor valuates crypto currency asset portfolios.

Usage:

    cryptor COMMAND [OPTION]...

Commands:

    init     create configuration directory and install example portfolios file
    valuate  calculate and display portfolio valuations
    help     display documentation

Options:

    -aggregate              Display combined portfolio valuations
    -confdir CONF_DIR       Specify directory containing data and cache files (default: $HOME/.cryptor)
    -currency CURRENCY      Display values in this fiat CURRENCY
    -date YYYY-MM-DD        Perform valuation using crypto prices as of date YYYY-MM-DD
    -portfolio PORTFOLIO    Process named portfolio (default: all portfolios)
    -force                  Unconditionally fetch crypto prices and exchange rates

Version:    v0.0.1 (linux/amd64)
Git commit: -
Built:      2022-12-22T19:26:45+13:00
Github:     https://github.com/srackham/cryptor
```

## Implementation and Usage Notes
- Crypto prices and exchange rates are cached locally to the cryptor configuration directory (default: `$HOME/.cryptor`). Price updates are only fetched when they are not found in the local cache files (unless the `-force` option is specified). Caching ensures minimal use of Web APIs which can be slow and are sometimes throttled.

- The `valuate` command values portfolio assets in the `portfolios.yaml` configuration file.

- Portfolio valuations are saved to the `$HOME/.cryptor/valuations.json` valuation history file.

- Valuations do not overwrite previously recorded valuations (this can be overridden with the `-force` option).

- Valuations of past dates (using the`-date` option) are made using historic crypto prices, otherwise today's crypto prices are used.

- All values are saved in USD (use the `-currency` option to display values in other fiat currencies).

- Values displayed in non-USD currencies are converted from USD values using today's exchange rates.

- If you specify the portfolio's `cost` (the portfolio's total cost in fiat currency) then portfolio gains and losses are calculated.

- The `cost` value is formatted like `<amount> <currency>`. The amount is mandatory; the currency is optional and defaults to `USD`; dollar and comma characters are ignored. Examples:

            $5,000.00 NZD     # Five thousand New Zealand dollars.
            1000              # One thousand US dollars

- Crypto and currency symbols are displayed in uppercase.
- Saved portfolio valuations are date stamped.
- Dates are specified and recorded as `YYYY-DD-MM` formatted strings.
- The `-portfolio` option can be specified multiple times.

- Crypto prices are fetched from [CoinGecko](https://www.coingecko.com/en/api); exchange rates are fetched from [exchangerate.host](https://exchangerate.host/).