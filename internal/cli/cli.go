package cli

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/logger"
	"github.com/srackham/cryptor/internal/portfolio"
	"github.com/srackham/cryptor/internal/price"
	"github.com/srackham/cryptor/internal/slice"
	"github.com/srackham/cryptor/internal/xrates"
)

// Build ldflags.
var (
	// VERS is the latest cryptor version tag. Set by linker -ldflags "-X main.VERS=..."
	VERS = "v0.0.1"
	// OS is the target operating system and architecture. Set by linker -ldflags "-X main.OS=..."
	OS = "-"
	// BUILD is the date the executable was built.
	BUILT = "-"
	// COMMIT is the Git commit hash.
	COMMIT = "-"
)

type cli struct {
	command         string
	executable      string
	configDir       string
	log             logger.Log
	portfolios      portfolio.Portfolios
	valuations      portfolio.Portfolios
	valuationsCache cache.Cache[portfolio.Portfolios]
	priceReader     price.PriceReader
	xrates          xrates.ExchangeRates
	opts            struct {
		aggregate  bool
		currency   string
		date       string
		force      bool
		portfolios []string
	}
}

// New creates a new cli context.
func New(api price.IPriceAPI) *cli {
	cli := cli{}
	cli.valuations = portfolio.Portfolios{}
	cli.valuationsCache = *cache.NewCache(&cli.valuations)
	cli.priceReader = price.NewPriceReader(api, &cli.log)
	cli.xrates = xrates.NewExchangeRates(&cli.log)
	return &cli
}

// Execute runs a command specified by CLI args.
func (cli *cli) Execute(args []string) error {
	var err error
	defer func() {
		if err != nil {
			cli.log.Error("%s", err.Error())
		}
	}()
	user, _ := user.Current()
	cli.configDir = filepath.Join(user.HomeDir, ".cryptor")
	cli.opts.currency = "USD"
	cli.opts.date = helpers.TodaysDate()
	err = cli.parseArgs(args)
	if err == nil {
		cli.priceReader.CacheFile = filepath.Join(cli.configDir, "crypto-prices.json")
		cli.xrates.CacheFile = filepath.Join(cli.configDir, "exchange-rates.json")
		cli.valuationsCache.CacheFile = filepath.Join(cli.configDir, "valuations.json")
		cli.priceReader.API.SetCacheDir(cli.configDir)
		switch cli.command {
		case "help":
			cli.help()
		case "init":
			err = cli.init()
		case "valuate":
			err = cli.valuate()
		default:
			err = fmt.Errorf("illegal command: " + cli.command)
		}
	}
	return err
}

// parseArgs parses and validate command-line arguments.
func (cli *cli) parseArgs(args []string) error {
	skip := false
	for i, opt := range args {
		if skip {
			skip = false
			continue
		}
		switch {
		case i == 0:
			cli.executable = opt
			if len(args) == 1 {
				cli.command = "help"
			}
		case i == 1:
			if opt == "-h" || opt == "--help" {
				opt = "help"
			}
			if !isCommand(opt) {
				return fmt.Errorf("illegal command: \"%s\"", opt)
			}
			cli.command = opt
		case opt == "-aggregate":
			cli.opts.aggregate = true
		case opt == "-force":
			cli.opts.force = true
		case slice.New("-confdir", "-currency", "-date", "-portfolio").Has(opt):
			// Process option argument.
			if i+1 >= len(args) {
				return fmt.Errorf("missing %s argument value", opt)
			}
			arg := args[i+1]
			switch opt {
			case "-confdir":
				cli.configDir = arg
			case "-currency":
				cli.opts.currency = strings.ToUpper(arg)
			case "-date":
				if !helpers.IsDateString(arg) {
					return fmt.Errorf("invalid date: \"%s\"", arg)
				}
				if strings.Compare(arg, helpers.TodaysDate()) == 1 {
					return fmt.Errorf("future date is not allowed: \"%s\"", arg)
				}
				cli.opts.date = arg
			case "-portfolio":
				cli.opts.portfolios = append(cli.opts.portfolios, arg)
			default:
				return fmt.Errorf("unexpected option: \"%s\"", opt)
			}
			skip = true
		default:
			return fmt.Errorf("illegal argument: \"%s\"", opt)
		}
	}
	if cli.command == "help" {
		return nil
	}
	return nil
}

// init implements the init command.
func (cli *cli) init() error {
	if !fsx.DirExists(cli.configDir) {
		cli.log.Highlight("creating configuration directory: \"%s\"", cli.configDir)
		if err := fsx.MkMissingDir(cli.configDir); err != nil {
			return err
		}
	}
	if fsx.FileExists(cli.configFile()) {
		return fmt.Errorf("portfolios file already exists: \"%s\"", cli.configFile())
	}
	cli.log.Highlight("installing example portfolios file: \"%s\"", cli.configFile())
	conf := `# Example cryptor portfolio configuration file

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
`
	if err := fsx.WriteFile(cli.configFile(), conf); err != nil {
		return err
	}
	return nil
}

// help implements the help command.
func (cli *cli) help() {
	github := "https://github.com/srackham/cryptor"
	summary := `Cryptor valuates crypto currency asset portfolios.

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

Version:    ` + VERS + " (" + OS + ")" + `
Git commit: ` + COMMIT + `
Built:      ` + BUILT + `
Github:     ` + github

	cli.log.Console("\n" + summary + "\n")
}

func isCommand(name string) bool {
	return slice.New("help", "init", "valuate").Has(name)
}

// plotValuations implements the `plot valuations` command.
// Plots the aggregate of the specified portfolios.
func (cli *cli) plotValuations() error {
	// TODO
	return nil
}

// plotAllocation implements the `plot allocation` command.
// Plots the aggregate of the specified portfolios.
func (cli *cli) plotAllocation() error {
	// TODO
	return nil
}

func (cli *cli) configFile() string {
	return filepath.Join(cli.configDir, "portfolios.yaml")
}

func (cli *cli) load() error {
	ps, err := portfolio.LoadPortfoliosFile(cli.configFile())
	if err != nil {
		return err
	}
	cli.portfolios = ps
	if err := cli.valuationsCache.Load(); err != nil {
		return err
	}
	if err := cli.xrates.Load(); err != nil {
		return err
	}
	if err := cli.priceReader.Load(); err != nil {
		return err
	}
	if err := cli.priceReader.API.LoadCacheFiles(); err != nil {
		return err
	}
	return nil
}

func (cli *cli) save() error {
	cli.valuations.SortByDateAndName()
	if err := cli.valuationsCache.Save(); err != nil {
		return err
	}
	if err := cli.xrates.Save(); err != nil {
		return err
	}
	if err := cli.priceReader.Save(); err != nil {
		return err
	}
	if err := cli.priceReader.API.SaveCacheFiles(); err != nil {
		return err
	}
	return nil
}

// valuate implements the valuate command.
func (cli *cli) valuate() error {
	date := cli.opts.date
	currency := cli.opts.currency
	if err := cli.load(); err != nil {
		return err
	}
	// Select portfolios.
	var ps portfolio.Portfolios
	if len(cli.opts.portfolios) > 0 {
		for _, name := range cli.opts.portfolios {
			i := cli.portfolios.FindByName(name)
			if i == -1 {
				return fmt.Errorf("missing portfolio: \"%s\"", name)
			}
			if ps.FindByName(name) != -1 {
				return fmt.Errorf("portfolio name can only be specified once: \"%s\"", name)
			}
			ps = append(ps, cli.portfolios[i])
		}
	} else {
		ps = cli.portfolios
	}
	// Evaluate portfolios.
	for i := range ps {
		ps[i].Date = date
		if err := ps[i].SetUSDValues(cli.priceReader, date, cli.opts.force); err != nil {
			return err
		}
		ps[i].SetAllocations()
		ps[i].Assets.SortByValue()
		if ps[i].Cost != "" {
			cost, err := cli.currencyToUSD(ps[i].Cost)
			if err != nil {
				return err
			}
			ps[i].USDCost = cost
		}
		// Update valuations history.
		if j := cli.valuations.FindByNameAndDate(ps[i].Name, ps[i].Date); j == -1 {
			cli.valuations = append(cli.valuations, ps[i])
		} else if date == helpers.TodaysDate() || cli.opts.force {
			cli.valuations[j] = ps[i]
		}
	}
	if cli.opts.aggregate {
		ps = []portfolio.Portfolio{ps.Aggregate("aggregate", date)}
	}
	// Print portfolios.
	xrate, err := cli.xrates.GetRate(cli.opts.currency, cli.opts.force)
	if err != nil {
		return err
	}
	ps.SortByDateAndName()
	cli.log.Console("")
	for _, p := range ps {
		s := fmt.Sprintf("NAME:  %s\nNOTES: %s\nDATE:  %s\nVALUE: %.2f %s",
			p.Name, p.Notes, p.Date, p.Value*xrate, currency)
		if p.USDCost > 0.00 {
			gains := p.Value - p.USDCost
			pcgains := gains / p.USDCost * 100
			s += fmt.Sprintf("\nCOST:  %.2f %s\nGAINS: %.2f (%.2f%%)", p.USDCost*xrate, currency, gains*xrate, pcgains)
		} else {
			s += "\nCOST:\nGAINS:"
		}
		if cli.opts.currency != "USD" {
			s += fmt.Sprintf("\nXRATE: 1 USD = %.2f %s", xrate, currency)
		} else {
			s += "\nXRATE:"
		}
		s += "\n            AMOUNT            VALUE   PERCENT            PRICE\n"
		for _, a := range p.Assets {
			value := a.Value * xrate
			s += fmt.Sprintf("%-5s %12.4f %12.2f %s    %5.2f%% %12.2f %s\n",
				a.Symbol,
				a.Amount,
				value,
				currency,
				a.Allocation,
				helpers.If(a.Amount > 0.0, value/a.Amount, 0),
				currency)
		}
		cli.log.Console("%s\n", s)
	}
	// Save valuations and cache files.
	if err := cli.save(); err != nil {
		return err
	}
	return nil
}

// currencyToUSD parses "<value>[<currency]" string and converts to USD.
func (cli *cli) currencyToUSD(cs string) (value float64, err error) {
	value, currency, err := portfolio.ParseCurrency(cs)
	if err != nil {
		return
	}
	rate, err := cli.xrates.GetRate(currency, cli.opts.force)
	if err != nil {
		return
	}
	if currency == "USD" && rate != 1.00 {
		err = fmt.Errorf("USD exchange rate should be 1.00: %f", rate)
		return
	}
	if rate == 0.00 {
		err = fmt.Errorf("exchange rate is zero: %s", currency)
		return
	}
	value = value / rate
	return
}
