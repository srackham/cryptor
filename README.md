# Cryptor

Cryptor valuates cryptocurrency asset portfolios.

-   Processes multiple asset portfolios.
-   Reports the current value and performance of cryptocurrency assets.

## Installation

If you have [Go](https://go.dev/) installed on your system then you can download and compile the latest version with this command:

    go install github.com/srackham/cryptor@latest

Pre-compiled binaries are also available on the [Cryptor releases page](https://github.com/srackham/cryptor/releases). Download the relevant release and extract the `cryptor` executable.

## Quick Start

1. Install an example configuration using the `cryptor init` command. For example:

    ```
    $ cryptor init
    creating configuration directory: "/home/srackham/.config/cryptor"
    installing example config file: "/home/srackham/.config/cryptor/config.yaml"
    installing example portfolios file: "/home/srackham/.config/cryptor/portfolios.yaml"
    creating cache directory: "/home/srackham/.cache/cryptor"
    creating data directory: "/home/srackham/.local/share/cryptor"
    ```

2. Edit the [`portfolios.yaml` configuration file](#portfolios-configuration-file) to match your portfolios.

3. Use the `cryptor valuate` command to calculate the current value the portfolios and their assets. For example:

    ```
    $ cryptor valuate

    NAME:  personal
    DATE:  2025-02-10
    TIME:  19:08:45
    VALUE: 55202.96 USD
    COST:  10000.00 USD
    GAINS: 45202.96 USD (452.03%)
                AMOUNT            VALUE    PERCENT       UNIT PRICE
    BTC         0.5000     48522.29 USD     87.90%     97044.58 USD
    ETH         2.5000      6580.68 USD     11.92%      2632.27 USD
    USDT      100.0000        99.99 USD      0.18%         1.00 USD
    ```

## Command Options

Run the `cryptor help` command to view all the commands and command options:

```
$ cryptor help

Usage:
    cryptor COMMAND [OPTION]...

Description:
    Cryptor valuates crypto currency asset portfolios.

Commands:
    init     create configuration directory and install default config and
             example portfolios files
    valuate  valuate, print and save portfolio valuations
    history  Print saved portfolio valuations
    help     display documentation

Options:
    -aggregate                  Include aggregated portfolios in printed valuation
    -aggregate-only             Only include aggregated portfolios in printed valuation
    -confdir CONF_DIR           Directory containing config, data and cache files
    -currency CURRENCY          Print fiat currency values denominated in CURRENCY
    -notes                      Include portfolio notes in the valuations
    -save                       Update the valuations file
    -portfolio PORTFOLIO        Print named portfolio valuation (default: all portfolios)
    -price SYMBOL=PRICE         Override the asset price of SYMBOL with PRICE (in USD)
    -format FORMAT              Set the valuate command output format ("json" or "yaml")

Config directory: /home/srackham/.config/cryptor
Cache directory:  /home/srackham/.cache/cryptor
Data directory:   /home/srackham/.local/share/cryptor

Version:    v0.3.0 (-)
Git commit: -
Built:      -
Github:     https://github.com/srackham/cryptor
```

## Implementation and Usage Notes

-   The `valuate` command values portfolio assets specified in the [`portfolios.yaml` configuration file](#portfolios-configuration-file).
-   All values are saved in USD (the `-currency` option can be used to display printed values in non-USD currencies).
-   Fiat currency exchange rates are cached locally and refreshed daily.
-   Asset and currency symbols are case insensitive and are converted to uppercase.
-   Saved portfolio valuations include the valuation's local `date` and `time`.
-   Dates are saved as `YYYY-DD-MM` formatted strings.
-   Times are saved as `hh:mm:ss` formatted strings.
-   Cryptocurrency prices are fetched using the [Binance HTTP ticker price](https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#symbol-price-ticker) API.
-   Fiat currency exchange rates are fetched using the [Open Exchange Rates](https://openexchangerates.org/) API.
-   By default valuations are printed in a human-friendly text format; use the `-format` option to print in JSON or YAML formats.
-   Currency values in JSON and YAML formats are always in USD.

-   The `-portfolio` option can be specified multiple times.
-   The `-price` option allows the user to override current asset prices in order to evaluate "what if" scenarios. Example:

        cryptor valuate -price btc=50000    # Price Bitcoin at $50,000 USD

-   The `-price` option can be specified multiple times.
-   The `-price` option could be used to include non-crypto assets, e.g. gold, in the portfolios.

-   Cryptor processes the following configuration, cache, and data files:

    -   `$HOME/.config/cryptor/config.yaml`: YAML formatted cryptor options
    -   `$HOME/.config/cryptor/portfolios.yaml`: YAML formatted portfolios
    -   `$HOME/.cache/cryptor/exchange-rates.json`: JSON formatted cached fiat currency exchange rates
    -   `$HOME/.local/share/data/cryptor/valuations.json`: JSON formatted valuations

-   Default locations for configuration, cache, and data files conform to the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/).
-   An alternate single directory for all files can be specified using the `-confdir` command option.
-   Use the `-confdir` option if you want to separate portfolios or groups of portfolios.

## Portfolios Configuration File

Asset portfolios are specified the [YAML](https://yaml.org/) formatted `portfolios.yaml` file in `$HOME/.config/cryptor/`). There are two portfolio file formats:

### Single-portfolio assets-only format
Contains a list of assets, one asset per line formatted like `<symbol>: <amount>`.
For example:

```yaml
BTC: 0.5
ETH: 2.5
USDT: 100
```

### Multi-portfolio format
Contains one or more portfolios each containing a list of assets along with optional portfolio name, notes and cost.

-   Portfolio names are unique and can only contain alphanumeric characters, underscores and dashes; the name `aggregate` is reserved.
-   If you specify a portfolio's `cost` amount (the total amount paid for the portfolio assets) then portfolio gains (or losses) are calculated.
-   The portfolio `cost` value is formatted like `<amount><symbol>`. The amount is mandatory; the currency symbol is optional and defaults to `USD`; dollar, comma and space characters are ignored; case insensitive. Examples:

        $5,000.00 NZD     # Five thousand New Zealand dollars.
        1000aud           # One thousand Australian dollars.
        .5                # Fifty cents USD.

Example multi-portfolios configuration file containing two portfolios:

```yaml
- name: personal
  notes: Personal portfolio notes.
  cost: $10,000.00 USD
  assets:
      BTC: 0.5
      ETH: 2.5
      USDT: 100

- name: business
  notes: |
      Business portfolio notes
      over multiple lines.
  cost: $20,000.00 USD
  assets:
      BTC: 1.0
```

## Valuations

-   If the `-save` option is specified the `valuate` command appends portfolio valuations to the `valuations.json` file located in the data configuration directory.
-   The `valuate` command prints and saves valuations in the same order that they occur in the portfolios configuration file.
-   The printed output can be customised using the `-portfolio`, `-aggregate` and `-aggregate-only` options.
-   The aggregate valuation is appended after the portfolio valuations with the portfolio name `aggregate`.
-   Saved valuations always include all portfolio valuations plus the aggregate valuation.
-   The `aggregate` portfolio is the aggregate of all portfolios, not just those specified by `-portfolio` options.
-   The `-portfolio`, `-aggregate` and `-aggregate-only` options apply to printed outputs.

## Post-processing Valuation Data

The [jq](https://github.com/jqlang/jq) command is useful for munging and extracting valuation data:

-   The following command pipes JSON output from the `valuate` command through a `jq` filter transforms it into CSV asset records:

          cryptor valuate -format json | jq -r '.[] | . as $p | .assets[] | [$p.name, $p.date,$p.time,.symbol,.amount,.value,.allocation] | @csv'

-   This command lists all saved portfolio valuations, includes a CSV header, and rounds numbers to two decimal places:

          cryptor history | jq -r '["NAME","DATE","TIME","VALUE","ROI"], (.[] | select(.cost>0) | [.name, .date, .time, (.value*100 | floor | ./100), ((.value-.cost)/.cost*100*100 | floor | ./100)]) | @csv'

-   The next command pipes the saved `personal` portfolio valuations through a `jq` filter to generate per-day CSV ROI (return on investment) records by selecting the first portfolio record of the day:

          cryptor history -portfolio personal | jq -r 'group_by(.date) | map(. | map(select(.cost > 0)) | sort_by(.time) | first) | .[] | select(.) | [.name, .date, .value, (.value-.cost)/.cost*100] | @csv'

## Plotting Portfolio Valuation Data

The cryptor repository includes an `examples` folder which contains bash scripts for plotting portfolio valuations. The scripts read the `cryptor` output on `stdin` and use [jq](https://stedolan.github.io/jq/) (to generate CSV plot data) and [gnuplot](http://www.gnuplot.info/) (to plot the CSV data).

### Portfolio valuation pie chart
The bash script `examples/plot-valuation.sh` plots a `cryptor` valuation. For example:

    cryptor valuate -portfolio personal -format json | examples/plot-valuation.sh

![Portfolio valuation pie chart](valuation-plot.png)

### Portfolio history chart
The bash script `examples/plot-history.sh` plots `cryptor` history data. For example:

    cryptor history -portfolio personal | examples/plot-history.sh

![Portfolio history chart](history-plot.png)
