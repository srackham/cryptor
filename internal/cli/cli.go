package cli

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/exchangerates"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/logger"
	"github.com/srackham/cryptor/internal/portfolio"
	"github.com/srackham/cryptor/internal/price"
	"github.com/srackham/cryptor/internal/slice"
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
	xrates          exchangerates.ExchangeRates
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
	cli.xrates = exchangerates.NewExchangeRates(&cli.log)
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
	cli.opts.date = helpers.DateNowString()
	err = cli.parseArgs(args)
	if err == nil {
		cli.priceReader.CacheFile = filepath.Join(cli.configDir, "crypto-prices.json")
		cli.xrates.CacheFile = filepath.Join(cli.configDir, "exchange-rates.json")
		cli.valuationsCache.CacheFile = filepath.Join(cli.configDir, "valuations.json")
		cli.priceReader.API.SetCacheDir(cli.configDir)
		switch cli.command {
		case "help":
			cli.help()
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
				return fmt.Errorf("illegal command: %q", opt)
			}
			cli.command = opt
		case opt == "-aggregate":
			cli.opts.aggregate = true
		case opt == "-force":
			cli.opts.force = true
		case opt == "-v":
			cli.log.Verbosity++
		case opt == "-vv":
			cli.log.Verbosity += 2
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
				cli.opts.currency = arg
			case "-date":
				if _, err := helpers.ParseDateString(arg, nil); err != nil {
					return fmt.Errorf("invalid date: %q", arg)
				}
				cli.opts.date = arg
			case "-portfolio":
				cli.opts.portfolios = append(cli.opts.portfolios, arg)
			default:
				return fmt.Errorf("unexpected option: %q", opt)
			}
			skip = true
		default:
			return fmt.Errorf("illegal argument: %q", opt)
		}
	}
	if cli.command == "help" {
		return nil
	}
	return nil
}

// help implements the help command.
func (cli *cli) help() {
	github := "https://github.com/srackham/cryptor"
	summary := `Cryptor reports crypto currency portfolio statistics.

Usage:

    cryptor COMMAND [OPTION]...

Commands:

    valuate    list portfolio valuations (default command)
    help        display documentation

Options:

    -aggregate              Combine portfolio valuations
    -confdir CONF_DIR       Directory containing data and cache files (default: $HOME/.cryptor)
    -currency CURRENCY      Display values in this CURRENCY
    -date DATE              Perform valuation using crypto prices as of DATE
    -portfolio PORTFOLIO    Process named portfolio (can be specified multiple times)
    -force                  Unconditionally fetch prices and exchange rates
    -v, -vv                 Increased verbosity

Version:    ` + VERS + " (" + OS + ")" + `
Git commit: ` + COMMIT + `
Built:      ` + BUILT + `
Github:     ` + github + ``

	cli.log.Console("\n" + summary)
}

func isCommand(name string) bool {
	return slice.New("help", "nop", "valuate").Has(name)
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

func (cli *cli) load() error {
	ps, err := portfolio.LoadPortfoliosFile(filepath.Join(cli.configDir, "portfolios.toml"))
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
	today := helpers.DateNowString()
	if err := cli.load(); err != nil {
		return err
	}
	var ps portfolio.Portfolios
	if len(cli.opts.portfolios) > 0 {
		for _, name := range cli.opts.portfolios {
			i := cli.portfolios.FindByName(name)
			if i == -1 {
				return fmt.Errorf("missing portfolio: %q", name)
			}
			if ps.FindByName(name) != -1 {
				return fmt.Errorf("portfolio name can only be specified once: %q", name)
			}
			ps = append(ps, cli.portfolios[i])
		}
	} else {
		ps = cli.portfolios
	}
	ps.SortByDateAndName()
	if cli.opts.aggregate {
		aggregate := ps.Aggregate("aggregate")
		ps = append(ps, aggregate)
	}
	prices, err := ps.GetPrices(cli.priceReader, date, cli.opts.force)
	if err != nil {
		return err
	}
	currency := strings.ToUpper(cli.opts.currency)
	xrate, err := cli.xrates.GetRate(currency, cli.opts.force && date == today)
	if err != nil {
		return err
	}
	cli.log.Console("")
	for _, p := range ps {
		p.SetUSDValues(prices)
		p.Date = date
		p.SetAllocations()
		p.Assets.SortByValue()
		if (p.Name != "aggregate" && !cli.opts.aggregate) || (p.Name == "aggregate" && cli.opts.aggregate) {
			s := fmt.Sprintf("NAME:  %s\nNOTES: %s\nDATE:  %s\nVALUE: %.2f %s",
				p.Name, p.Notes, p.Date, p.Value*xrate, currency)
			if p.Cost != "" {
				cost, err := cli.toUSD(p.Cost)
				if err != nil {
					return err
				}
				gains := p.Value - cost
				pcgains := helpers.If(cost != 0.00, gains/cost*100, 0)
				s += fmt.Sprintf("\nCOST:  %.2f %s\nGAINS: %.2f (%.2f%%)", cost*xrate, currency, gains*xrate, pcgains)
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
		if p.Name != "aggregate" {
			i := cli.valuations.FindByNameAndDate(p.Name, p.Date)
			if i == -1 {
				cli.valuations = append(cli.valuations, p)
			} else {
				(cli.valuations)[i] = p
			}
		}
	}
	if err := cli.save(); err != nil {
		return err
	}
	return nil
}

func (cli *cli) toUSD(s string) (value float64, err error) {
	value, currency, err := portfolio.ParseCurrency(s)
	if err != nil {
		return
	}
	rate, err := cli.xrates.GetRate(currency, false)
	if err != nil {
		return
	}
	if rate == 0.00 {
		err = fmt.Errorf("exchange rate is zero: %s", currency)
		return
	}
	value = value / rate
	return
}
