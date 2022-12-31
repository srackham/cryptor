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
$ cryptor

Cryptor valuates crypto currency asset portfolios.

Usage:

    cryptor COMMAND [OPTION]...

Commands:

    init     create configuration directory and install example portfolios file
    valuate  calculate and display portfolio valuations
    history  display saved portfolio valuations from the valuations history
    help     display documentation

Options:

    -aggregate              Display portfolio valuations aggregated by date
    -confdir CONF_DIR       Directory containing data and cache files (default: $HOME/.cryptor)
    -currency CURRENCY      Display values in this fiat CURRENCY
    -date DATE              Valuation date, YYYY-MM-DD format or integer day offset: 0,-1,-2...
    -format FORMAT          Print format: text, json
    -portfolio PORTFOLIO    Process named portfolio (default: all portfolios)
    -force                  Unconditionally fetch crypto prices and exchange rates

Version:    v0.1.0 (linux/amd64)
Git commit: -
Built:      2022-12-31T15:02:13+13:00
Github:     https://github.com/srackham/cryptor
```

## Implementation and Usage Notes
- Crypto prices and exchange rates are cached locally to the cryptor configuration directory (default: `$HOME/.cryptor`). Price updates are only fetched when they are not found in the local cache files (unless the `-force` option is specified). Caching ensures minimal use of Web APIs which can be slow and are sometimes throttled.

- The `valuate` command values portfolio assets in the `portfolios.yaml` configuration file.

- Portfolio valuations are saved to the `$HOME/.cryptor/valuations.json` valuation history file.

- Valuations do not overwrite previously recorded valuations (this can be overridden with the `-force` option).

- Valuations of past dates (using the`-date` option) are made using historic crypto prices, otherwise today's crypto prices are used.

- The `history` command displays previously saved valuations from the valuations history file, if no matching valuations are found nothing will be displayed.

- All values are saved in USD (use the `-currency` option to display values in other fiat currencies).

- Values displayed in non-USD currencies are converted from USD values using today's exchange rates.

- If you specify the portfolio's `cost` (the portfolio's total cost in fiat currency) then portfolio gains and losses are calculated.

- The `cost` value is formatted like `<amount> <currency>`. The amount is mandatory; the currency is optional and defaults to `USD`; dollar and comma characters are ignored. Examples:

            $5,000.00 NZD     # Five thousand New Zealand dollars.
            1000              # One thousand US dollars

- Crypto and currency symbols are displayed in uppercase.
- Saved portfolio valuations are date stamped.
- Dates are specified either as `YYYY-DD-MM` formatted strings or as an integer date offset: `0` for today, `-1` for yesterday etc. For example `-date -7` specifies the date one week ago.
- Dates are recorded as `YYYY-DD-MM` formatted strings.
- The `-portfolio` option can be specified multiple times.

- Crypto prices are fetched from [CoinGecko](https://www.coingecko.com/en/api); exchange rates are fetched from [exchangerate.host](https://exchangerate.host/).


## Portfolio Valuations Data
There are two formats for printing portfolio valuations:

- `text`: human-readable text format.
- `json`: JSON format.

Other formats such as CSV can be extracted from the JSON formatted data using external tools. One such tool is [jq](https://stedolan.github.io/jq/). Here are some examples of `jq` filters:

### CSV assets history
The following command pipes JSON valuations history through a `jq` filter transforming it into CSV asset records: `<name>, <date>, <symbol>, <amount>, <value>, <allocation percentage>` records:

```
cryptor history -format json | jq -r '.[] | . as $p | .assets[] | [$p.name, $p.date,.symbol,.amount,.value,.allocation] | @csv'
```
Output:

```
"joint","2022-12-24","BTC",0.5,8414.9052734375,73.43150301510853
"joint","2022-12-24","ETH",2.5,3044.6249389648438,26.56849698489147
"personal","2022-12-24","BTC",0.5,8414.9052734375,72.79373770058
"personal","2022-12-24","ETH",2.5,3044.6249389648438,26.33774498962544
"personal","2022-12-24","USDC",100,100.39999485015869,0.8685173097945648
"portfolio1","2022-12-24","BTC",0.25,4207.45263671875,100
```

### CSV portfolio ROI history
The following command pipes JSON valuations history through a `jq` filter transforming costed valuations (valuations with `usdcost>0`) into CSV portfolio ROI (return on investment) records: `<name>, <date>, <value>, <percent ROI>` records:

```
cryptor history -format json | jq -r '.[] | select(.usdcost>0) | [.name, .date, .value, (.value-.usdcost)/.usdcost*100] | @csv'
```
Output:

```
"personal","2022-11-01",14278.401210995735,126.2568306455474
"personal","2022-12-24",11559.930207252502,83.40291671014384
```

The same query with a CSV header row and numbers rounded to two decimal places:

```
$ cryptor history -format json | jq -r '["NAME","DATE","VALUE","ROI"], (.[] | select(.usdcost>0) | [.name, .date, (.value*100 | floor | ./100), ((.value-.usdcost)/.usdcost*100*100 | floor | ./100)]) | @csv'
```
Output:

```
"NAME","DATE","VALUE","ROI"
"personal","2022-11-01",14278.4,126.25
"personal","2022-12-24",11559.93,83.4
```

## Plotting Portfolio Valuation Data
The cryptor repository includes an `examples` folder which contains bash scripts for plotting portfolio valuations. These scripts:

- Use [jq](https://stedolan.github.io/jq/) (to extract CSV plot data from JSON) and [gnuplot](http://www.gnuplot.info/) (to plot the data).
- Consume JSON formatted `cryptor` portfolio valuations from `stdin`.
- Accept an optional chart title as first script argument.

### Portfolio valuation pie chart
The bash script `examples/plot-valuation.sh` plots a `cryptor` valuation. For example:

    cryptor valuate -format json -portfolio personal | examples/plot-valuation.sh

![Portfolio valuation pie chart](valuation-plot.png)