package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/mock"
	"github.com/srackham/cryptor/internal/portfolio"
)

func mockCli(t *testing.T) *cli {
	ctx := mock.NewContext()
	if !fsx.DirExists(ctx.CacheDir) {
		if err := os.Mkdir(ctx.CacheDir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return New(&ctx)
}

func TestLoadConfig(t *testing.T) {
	cli := mockCli(t)
	portfoliosFile := path.Join(cli.ConfigDir, "portfolios.yaml")
	ps, err := cli.loadConfigFile(portfoliosFile)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 3, len(ps))
	assert.Equal(t, 3, len(ps[0].Assets))
	assert.Equal(t, "personal", ps[0].Name)
	i := ps[0].Assets.Find("BTC")
	assert.PassIf(t, i != -1, "missing asset: BTC")
	assert.Equal(t, portfolio.Asset{
		Symbol: "BTC",
		Amount: 0.5,
		Value:  0.0,
	}, ps[0].Assets[i])
	i = ps[1].Assets.Find("ETH")
	assert.PassIf(t, i != -1, "missing asset: ETH")
	assert.Equal(t, 2, len(ps[1].Assets))
	assert.Equal(t, "joint", ps[1].Name)
	assert.Equal(t, portfolio.Asset{
		Symbol: "ETH",
		Amount: 2.5,
		Value:  0.0,
	}, ps[1].Assets[i])
}

func TestValuate(t *testing.T) {
	cli := mockCli(t)
	portfoliosFile := path.Join(cli.ConfigDir, "portfolios.yaml")
	ps, err := cli.loadConfigFile(portfoliosFile)
	assert.PassIf(t, err == nil, "error reading portfolios file")
	p := ps[0]
	reader := cli.priceReader
	err = p.SetUSDValues(reader)
	assert.Equal(t, 52600.0, p.Value)
	assert.PassIf(t, err == nil, "error pricing portfolio: %v", err)
	p.Assets.Sort()
	assert.Equal(t, 50000.0, p.Assets[0].Value)
	assert.Equal(t, 2500.0, p.Assets[1].Value)
	assert.Equal(t, 100.0, p.Assets[2].Value)
	p.Assets[0].Value = 1000.00
	p.Assets.Sort()
	assert.Equal(t, 2500.0, p.Assets[0].Value)
	assert.Equal(t, 1000.0, p.Assets[1].Value)
	assert.Equal(t, 100.0, p.Assets[2].Value)
}

func TestParseArgs(t *testing.T) {
	var cli *cli
	var err error
	parse := func(cmd string) {
		args := strings.Split(cmd, " ")
		cli = mockCli(t)
		err = cli.parseArgs(args)
	}
	parse("cryptor")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, "help", cli.command)
	parse("cryptor invalid-command")
	assert.Equal(t, `invalid command: "invalid-command"`, err.Error())
}

func exec(cli *cli, cmd string) (string, string, error) {
	args := strings.Split(cmd, " ")
	err := cli.Execute(args...)
	return cli.Stdout.(*bytes.Buffer).String(), cli.Stderr.(*bytes.Buffer).String(), err
}

func TestValuateCmd(t *testing.T) {
	cli := mockCli(t)
	stdout, _, err := exec(cli, "cryptor valuate -no-save")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 3, len(cli.valuation))
	assert.PassIf(t, cli.valuation.FindByNameAndDate("personal", "2000-12-01") != -1, "missing valuation")
	wanted := `
NAME:  personal
NOTES:
DATE:  2000-12-01
TIME:  12:30:00
VALUE: 52600.00 USD
COST:  6666.67 USD
GAINS: 45933.33 USD (689.00%)
XRATE:
            AMOUNT            VALUE    PERCENT            PRICE
BTC         0.5000     50000.00 USD     95.06%    100000.00 USD
ETH         2.5000      2500.00 USD      4.75%      1000.00 USD
USDC      100.0000       100.00 USD      0.19%         1.00 USD

NAME:  joint
NOTES:
DATE:  2000-12-01
TIME:  12:30:00
VALUE: 52500.00 USD
COST:
GAINS:
XRATE:
            AMOUNT            VALUE    PERCENT            PRICE
BTC         0.5000     50000.00 USD     95.24%    100000.00 USD
ETH         2.5000      2500.00 USD      4.76%      1000.00 USD

NAME:  portfolio1
NOTES:
DATE:  2000-12-01
TIME:  12:30:00
VALUE: 25000.00 USD
COST:
GAINS:
XRATE:
            AMOUNT            VALUE    PERCENT            PRICE
BTC         0.2500     25000.00 USD    100.00%    100000.00 USD
`
	wanted = helpers.StripTrailingSpaces(wanted)
	stdout = helpers.StripTrailingSpaces(stdout)
	assert.EqualStrings(t, wanted, stdout)

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor valuate -no-save -notes")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 3, len(cli.valuation))
	assert.PassIf(t, cli.valuation.FindByNameAndDate("personal", "2000-12-01") != -1, "missing valuation")
	wanted = `
NAME:  personal
NOTES: ## Personal Portfolio
- 7-Jan-2023: Migrated to new h/w wallet.

DATE:  2000-12-01
TIME:  12:30:00
VALUE: 52600.00 USD
COST:  6666.67 USD
GAINS: 45933.33 USD (689.00%)
XRATE:
            AMOUNT            VALUE    PERCENT            PRICE
BTC         0.5000     50000.00 USD     95.06%    100000.00 USD
ETH         2.5000      2500.00 USD      4.75%      1000.00 USD
USDC      100.0000       100.00 USD      0.19%         1.00 USD

NAME:  joint
NOTES: Joint Portfolio
DATE:  2000-12-01
TIME:  12:30:00
VALUE: 52500.00 USD
COST:
GAINS:
XRATE:
            AMOUNT            VALUE    PERCENT            PRICE
BTC         0.5000     50000.00 USD     95.24%    100000.00 USD
ETH         2.5000      2500.00 USD      4.76%      1000.00 USD

NAME:  portfolio1
NOTES:
DATE:  2000-12-01
TIME:  12:30:00
VALUE: 25000.00 USD
COST:
GAINS:
XRATE:
            AMOUNT            VALUE    PERCENT            PRICE
BTC         0.2500     25000.00 USD    100.00%    100000.00 USD
`
	wanted = helpers.StripTrailingSpaces(wanted)
	stdout = helpers.StripTrailingSpaces(stdout)
	assert.EqualStrings(t, wanted, stdout)

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor valuate -no-save -aggregate-only -notes")
	assert.PassIf(t, err == nil, "%v", err)
	wanted = `
NAME:  aggregate
NOTES: joint, personal, portfolio1
DATE:  2000-12-01
TIME:  12:30:00
VALUE: 130100.00 USD
COST:
GAINS:
XRATE:
            AMOUNT            VALUE    PERCENT            PRICE
BTC         1.2500    125000.00 USD     96.08%    100000.00 USD
ETH         5.0000      5000.00 USD      3.84%      1000.00 USD
USDC      100.0000       100.00 USD      0.08%         1.00 USD
`
	wanted = helpers.StripTrailingSpaces(wanted)
	stdout = helpers.StripTrailingSpaces(stdout)
	assert.EqualStrings(t, wanted, stdout)

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor valuate -no-save -portfolio portfolio1 -portfolio personal")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 3, len(cli.valuation))
	assert.PassIf(t, cli.valuation.FindByNameAndDate("personal", "2000-12-01") != -1, "missing valuation")
	assert.PassIf(t, cli.valuation.FindByNameAndDate("portfolio1", "2000-12-01") != -1, "missing valuation")

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor valuate -no-save -aggregate-only -format json")
	assert.PassIf(t, err == nil, "%v", err)
	wanted = `
[
  {
    "name": "aggregate",
    "notes": "joint, personal, portfolio1",
    "date": "2000-12-01",
    "time": "12:30:00",
    "value": 130100,
    "cost": 0,
    "assets": [
      {
        "symbol": "BTC",
        "amount": 1.25,
        "value": 125000,
        "allocation": 96.07993850883936
      },
      {
        "symbol": "ETH",
        "amount": 5,
        "value": 5000,
        "allocation": 3.843197540353574
      },
      {
        "symbol": "USDC",
        "amount": 100,
        "value": 100,
        "allocation": 0.07686395080707148
      }
    ]
  }
]
`
	assert.EqualStrings(t, wanted, stdout)

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor valuate -no-save -aggregate-only -format yaml")
	assert.PassIf(t, err == nil, "%v", err)
	wanted = `
- name: aggregate
  notes: joint, personal, portfolio1
  date: "2000-12-01"
  time: "12:30:00"
  value: 130100
  cost: 0
  assets:
    - symbol: BTC
      amount: 1.25
      value: 125000
      allocation: 96.07993850883936
    - symbol: ETH
      amount: 5
      value: 5000
      allocation: 3.843197540353574
    - symbol: USDC
      amount: 100
      value: 100
      allocation: 0.07686395080707148
`
	assert.EqualStrings(t, wanted, stdout)

	cli = mockCli(t)
	if fsx.FileExists(cli.xrates.CacheFile()) {
		err = os.Remove(cli.xrates.CacheFile())
		assert.PassIf(t, err == nil, "%v", err)
	}
	stdout, _, err = exec(cli, "cryptor valuate -no-save -aggregate-only -currency NZD")
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, fsx.FileExists(cli.xrates.CacheFile()), "missing exchange rates cache file: \"%v\"", cli.xrates.CacheFile())
	got, err := fsx.ReadFile(cli.xrates.CacheFile())
	assert.PassIf(t, err == nil, "%v", err)
	wanted = `{
  "2000-12-01": {
    "AUD": 1.6,
    "NZD": 1.5,
    "USD": 1
  }
}`
	assert.EqualStrings(t, wanted, got)
	wanted = `
NAME:  aggregate
NOTES: joint, personal, portfolio1
DATE:  2000-12-01
TIME:  12:30:00
VALUE: 195150.00 NZD
COST:
GAINS:
XRATE: 1 USD = 1.50 NZD
            AMOUNT            VALUE    PERCENT            PRICE
BTC         1.2500    187500.00 NZD     96.08%    150000.00 NZD
ETH         5.0000      7500.00 NZD      3.84%      1500.00 NZD
USDC      100.0000       150.00 NZD      0.08%         1.50 NZD
`
	wanted = helpers.StripTrailingSpaces(wanted)
	stdout = helpers.StripTrailingSpaces(stdout)
	assert.EqualStrings(t, wanted, stdout)
}

func TestHelpCmd(t *testing.T) {
	cli := mockCli(t)
	stdout, _, err := exec(cli, "cryptor help")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, stdout, `
Usage:
    cryptor COMMAND [OPTION]...

Description:
    Cryptor valuates crypto currency asset portfolios.

Commands:
    init     create configuration directory and install default config and
             example portfolios files
    valuate  valuate, print and save portfolio valuations
    history  Print saved portfolio valuations
    help     display documentation`)
}

func TestInitCmd(t *testing.T) {
	tmpdir := mock.MkdirTemp(t)
	cli := mockCli(t)
	err := cli.Execute("cryptor", "init", "-confdir", tmpdir)
	assert.PassIf(t, err == nil, "%v", err)
	stdout := cli.Stdout.(*bytes.Buffer).String()
	assert.Contains(t, stdout, `installing example config file:`)
	s, err := fsx.ReadFile(cli.configFile())
	assert.PassIf(t, err == nil, "%v", err)
	assert.EqualStrings(t, `# NOTE: The xrates-appid option is only necessary if the -currency command option is used.
# Open Exchange Rates App ID (https://openexchangerates.org/)
xrates-appid: YOUR_APP_ID`, s)
	assert.Contains(t, stdout, `installing example portfolios file:`)
	s, err = fsx.ReadFile(cli.portfoliosFile())
	assert.PassIf(t, err == nil, "%v", err)
	assert.EqualStrings(t, `# Example cryptor portfolios configuration file containing two portfolios installed by 'cryptor init' command.

- name:  personal
  notes: Personal portfolio notes.
  cost: $10,000.00 USD
  assets:
    BTC: 0.5
    ETH: 2.5
    USDC: 100

- name:  business
  notes: |
    Business portfolio notes
    over multiple lines.
  cost: $20,000.00 USD
  assets:
    BTC: 1.0`, s)
	err = cli.Execute("cryptor", "init", "-confdir", tmpdir)
	assert.PassIf(t, err == nil, "%v", err)
	stdout = cli.Stdout.(*bytes.Buffer).String()
	assert.Contains(t, stdout, `config file already exists:`)
	assert.Contains(t, stdout, `portfolios file already exists:`)
}

func TestSaveValuation(t *testing.T) {
	tmpdir := mock.MkdirTemp(t)

	cli := mockCli(t)
	cli.DataDir = tmpdir
	_, _, err := exec(cli, "cryptor valuate")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 3, len(cli.valuation))
	savedValuations := portfolio.Portfolios{}
	savedValuations = append(savedValuations, cli.valuation...)
	savedValuations = append(savedValuations, cli.aggregate)
	valuationsFile := cli.valuationsFile("json")
	loadedValuations, err := portfolio.LoadValuations(valuationsFile)
	assert.PassIf(t, err == nil, "error loading valuations file: \"%v\": %v", valuationsFile, err)
	assert.Equal(t, "personal", loadedValuations[0].Name)
	assert.Equal(t, "joint", loadedValuations[1].Name)
	assert.Equal(t, "portfolio1", loadedValuations[2].Name)
	assert.Equal(t, "aggregate", loadedValuations[3].Name)
	assert.PassIf(t, len(savedValuations) == 4, "valuations file: \"%v\": expected 4 portfolios, got %d", valuationsFile, len(savedValuations))
	assert.PassIf(t, reflect.DeepEqual(savedValuations, loadedValuations),
		"valuations file: \"%v\": expected:\n%v\n\ngot:\n%v", valuationsFile, savedValuations, loadedValuations)

	cli = mockCli(t)
	cli.DataDir = tmpdir
	_, _, err = exec(cli, "cryptor valuate")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 3, len(cli.valuation))
	savedValuations = append(savedValuations, cli.valuation...)
	savedValuations = append(savedValuations, cli.aggregate)
	loadedValuations, err = portfolio.LoadValuations(valuationsFile)
	assert.PassIf(t, err == nil, "error loading valuations file: \"%v\": %v", valuationsFile, err)
	assert.Equal(t, "personal", loadedValuations[4].Name)
	assert.Equal(t, "joint", loadedValuations[5].Name)
	assert.Equal(t, "portfolio1", loadedValuations[6].Name)
	assert.Equal(t, "aggregate", loadedValuations[7].Name)
	assert.PassIf(t, len(savedValuations) == 8, "valuations file: \"%v\": expected 8 portfolios, got %d", valuationsFile, len(savedValuations))
	assert.PassIf(t, reflect.DeepEqual(savedValuations, loadedValuations),
		"valuations file: \"%v\": expected:\n%v\n\ngot:\n%v", valuationsFile, savedValuations, loadedValuations)
}

func TestMissingPortfolio(t *testing.T) {
	cli := mockCli(t)
	_, stderr, err := exec(cli, "cryptor valuate -no-save -portfolio non-existent")
	assert.FailIf(t, err == nil, "non-existent portfolio should generate an error")
	assert.Contains(t, stderr, "missing portfolio: \"non-existent\"")
}

func TestHistoryCmd(t *testing.T) {
	cli := mockCli(t)
	stdout, _, err := exec(cli, "cryptor history")
	assert.PassIf(t, err == nil, "%v", err)
	valuations := portfolio.Portfolios{}
	err = json.Unmarshal([]byte(stdout), &valuations)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 15, len(valuations))
	valuationsFile := cli.valuationsFile("json")
	savedValuations, err := portfolio.LoadValuations(valuationsFile)
	assert.PassIf(t, reflect.DeepEqual(savedValuations, valuations),
		"valuations file: \"%v\": expected:\n%v\n\ngot:\n%v", valuationsFile, savedValuations, valuations)

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor history -portfolio joint")
	assert.PassIf(t, err == nil, "%v", err)
	valuations = portfolio.Portfolios{}
	err = json.Unmarshal([]byte(stdout), &valuations)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 7, len(valuations))
	assert.Equal(t, "joint", valuations[0].Name)

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor history -portfolio joint -portfolio personal")
	assert.PassIf(t, err == nil, "%v", err)
	valuations = portfolio.Portfolios{}
	err = json.Unmarshal([]byte(stdout), &valuations)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 14, len(valuations))

	cli = mockCli(t)
	stdout, _, err = exec(cli, "cryptor history -portfolio aggregate")
	assert.PassIf(t, err == nil, "%v", err)
	valuations = portfolio.Portfolios{}
	err = json.Unmarshal([]byte(stdout), &valuations)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1, len(valuations))
	assert.Equal(t, "aggregate", valuations[0].Name)

	cli = mockCli(t)
	_, stderr, err := exec(cli, "cryptor history -portfolio non-existent")
	assert.FailIf(t, err == nil, "non-existent portfolio should generate an error")
	assert.Contains(t, stderr, "no valuations found")
}

func TestNoConfigFile(t *testing.T) {
	tmpdir := mock.MkdirTemp(t)
	cli := mockCli(t)
	err := fsx.WriteFile(path.Join(tmpdir, "portfolios.yaml"), `# Minimal portfolio
- assets:
      BTC: 0.25`)
	assert.PassIf(t, err == nil, "%v", err)
	cli.ConfigDir = tmpdir // No config.yaml file
	cli.DataDir = tmpdir
	_, _, err = exec(cli, "cryptor valuate")
	assert.PassIf(t, err == nil, "missing config.yaml file should not generate an error: %v", err)
}

func TestParsePriceOption(t *testing.T) {
	testCases := []struct {
		name           string
		priceOption    string
		expectedSymbol string
		expectedPrice  float64
		expectedErr    bool
	}{
		{
			name:           "Valid price option",
			priceOption:    "BTC=12345.67",
			expectedSymbol: "BTC",
			expectedPrice:  12345.67,
			expectedErr:    false,
		},
		{
			name:           "Valid price option with spaces",
			priceOption:    "  ETH  =  123.45  ",
			expectedSymbol: "ETH",
			expectedPrice:  123.45,
			expectedErr:    false,
		},
		{
			name:           "Invalid price option - missing equals",
			priceOption:    "BTC12345.67",
			expectedSymbol: "",
			expectedPrice:  0,
			expectedErr:    true,
		},
		{
			name:           "Invalid price option - multiple equals",
			priceOption:    "BTC=12345.67=890",
			expectedSymbol: "",
			expectedPrice:  0,
			expectedErr:    true,
		},
		{
			name:           "Invalid price option - invalid symbol",
			priceOption:    "BTC@=12345.67",
			expectedSymbol: "",
			expectedPrice:  0,
			expectedErr:    true,
		},
		{
			name:           "Invalid price option - invalid price",
			priceOption:    "BTC=abc",
			expectedSymbol: "",
			expectedPrice:  0,
			expectedErr:    true,
		},
		{
			name:           "Valid price option - negative price",
			priceOption:    "BTC=-123.45",
			expectedSymbol: "BTC",
			expectedPrice:  -123.45,
			expectedErr:    false,
		},
		{
			name:           "Valid price option - zero price",
			priceOption:    "BTC=0",
			expectedSymbol: "BTC",
			expectedPrice:  0,
			expectedErr:    false,
		},
		{
			name:           "Valid price option - symbol with hyphen and underscore",
			priceOption:    "BTC-USD_PERP=123.45",
			expectedSymbol: "BTC-USD_PERP",
			expectedPrice:  123.45,
			expectedErr:    false,
		},
		{
			name:           "Valid price option - symbol with numbers",
			priceOption:    "BTC123=123.45",
			expectedSymbol: "BTC123",
			expectedPrice:  123.45,
			expectedErr:    false,
		},
		{
			name:           "Empty price option",
			priceOption:    "",
			expectedSymbol: "",
			expectedPrice:  0,
			expectedErr:    true,
		},
		{
			name:           "Price option with only symbol",
			priceOption:    "BTC=",
			expectedSymbol: "",
			expectedPrice:  0,
			expectedErr:    true,
		},
		{
			name:           "Price option with only price",
			priceOption:    "=123.45",
			expectedSymbol: "",
			expectedPrice:  0,
			expectedErr:    true,
		},
	}
	for _, tc := range testCases {
		symbol, price, err := ParsePriceOption(tc.priceOption)
		if tc.expectedErr {
			if err == nil {
				t.Errorf("%s: expected error, but got nil", tc.name)
			}
		} else {
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tc.name, err)
			}
			if symbol != tc.expectedSymbol {
				t.Errorf("%s: expected symbol %q, but got %q", tc.name, tc.expectedSymbol, symbol)
			}
			if price != tc.expectedPrice {
				t.Errorf("%s: expected price %f, but got %f", tc.name, tc.expectedPrice, price)
			}
		}
	}
}
