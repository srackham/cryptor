package cli

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/srackham/cryptor/internal/cache"
	"github.com/srackham/cryptor/internal/exchangerates"
	"github.com/srackham/cryptor/internal/fsx"
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
	command      string
	executable   string
	configDir    string
	configFile   string
	log          logger.Log
	portfolios   portfolio.Portfolios
	history      portfolio.Portfolios
	historyCache cache.Cache[portfolio.Portfolios]
	priceReader  price.PriceReader
	xrates       exchangerates.ExchangeRates
	opts         struct {
		aggregate  bool     // If true then combine portfolios
		currency   string   // Symbol of denominated fiat currency (defaults to USD).
		date       string   // Use previously recorded valuate from history file.
		refresh    bool     // If true unconditionally update prices and exchange rates.
		portfolios []string // Portfolios to process (default to all portfolios)
	}
}

// New creates a new cli context.
func New(api price.IPriceAPI) *cli {
	cli := cli{}
	cli.history = portfolio.Portfolios{}
	cli.historyCache = *cache.NewCache(&cli.history)
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
	cli.configFile = filepath.Join(cli.configDir, "portfolios.toml")
	cli.opts.currency = "USD"
	cli.opts.date = helpers.DateNowString()
	err = cli.parseArgs(args)
	if err == nil {
		cli.priceReader.CacheFile = filepath.Join(cli.configDir, "crypto-prices.json")
		cli.xrates.CacheFile = filepath.Join(cli.configDir, "exchange-rates.json")
		cli.historyCache.CacheFile = filepath.Join(cli.configDir, "history.json")
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
		case opt == "-refresh":
			cli.opts.refresh = true
		case opt == "-v":
			cli.log.Verbosity++
		case opt == "-vv":
			cli.log.Verbosity += 2
		case slice.New("-conf", "-confdir", "-currency", "-date", "-portfolio").Has(opt):
			// Process option argument.
			if i+1 >= len(args) {
				return fmt.Errorf("missing %s argument value", opt)
			}
			arg := args[i+1]
			switch opt {
			case "-conf":
				cli.configFile = arg
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
    -conf CONF              Configuration file (default: CONF_DIR/portfolios.toml)
    -confdir CONF_DIR       Directory containing data and cache files (default: $HOME/.cryptor)
    -currency CURRENCY      Display values in this CURRENCY
    -date DATE              Perform valuation using crypto prices as of DATE
    -portfolio PORTFOLIO    Process named portfolio (can be specified multiple times)
    -refresh                Fetch the latest prices and exchange rates
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

// plotHistory implements the `plot history` command.
// Plots the aggregate of the specified portfolios.
func (cli *cli) plotHistory() error {
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
	var err error
	if err = cli.loadPortfolios(); err != nil {
		return err
	}
	if err := cli.historyCache.LoadCacheFile(); err != nil {
		return err
	}
	if err := cli.xrates.LoadCacheFile(); err != nil {
		return err
	}
	if err := cli.priceReader.LoadCacheFile(); err != nil {
		return err
	}
	if err := cli.priceReader.API.LoadCacheFiles(); err != nil {
		return err
	}
	return nil
}

func (cli *cli) save() error {
	if err := cli.historyCache.SaveCacheFile(); err != nil {
		return err
	}
	for k, _ := range *cli.xrates.CacheData {
		// We only use the current exchange rates so delete non-current entries.
		if k != helpers.DateNowString() {
			delete(*cli.xrates.CacheData, k)
		}
	}
	if err := cli.xrates.SaveCacheFile(); err != nil {
		return err
	}
	if err := cli.priceReader.SaveCacheFile(); err != nil {
		return err
	}
	if err := cli.priceReader.API.SaveCacheFiles(); err != nil {
		return err
	}
	return nil
}

// valuate implements the valuate command.
func (cli *cli) valuate() error {
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
	if cli.opts.aggregate {
		aggregate := ps.Aggregate("__aggregate__", "Combined Portfolios")
		ps = append(ps, aggregate)
	}
	prices, err := ps.GetPrices(cli.priceReader, cli.opts.date, cli.opts.refresh)
	if err != nil {
		return err
	}
	currency := strings.ToUpper(cli.opts.currency)
	xrate, err := cli.xrates.GetRate(currency, helpers.DateNowString(), cli.opts.refresh) // Use current exchange rates.
	if err != nil {
		return err
	}
	cli.log.Console("")
	for _, p := range ps {
		p.SetUSDValues(prices)
		p.SetTimeStamp(cli.opts.date, cli.opts.refresh)
		p.SetAllocations()
		p.Assets.SortByValue()
		if p.Name != "__aggregate__" && !cli.opts.aggregate || p.Name == "__aggregate__" {
			s := fmt.Sprintf(`Name:        %s
Description: %s
TimeStamp:   %s %s
Value:       %.2f %s

`,
				p.Name, p.Description, p.Date, p.Time, p.USD*xrate, currency)
			s += "            AMOUNT            PRICE       UNIT PRICE\n"
			for _, a := range p.Assets {
				value := a.USD * xrate
				s += fmt.Sprintf("%-5s %12.4f %12.2f %s %12.2f %s    %5.2f%%\n",
					a.Symbol,
					a.Amount,
					value,
					currency,
					helpers.If(a.Amount > 0.0, value/a.Amount, 0),
					currency,
					a.Allocation)
			}
			cli.log.Console("%s\n", s)
		}
		// Only update history with today's valuation.
		if p.Name != "__aggregate__" && cli.opts.date == helpers.DateNowString() {
			cli.history.UpdateHistory(p)
		}
	}
	if err := cli.save(); err != nil { // Don't update unless the command succeeds.
		return err
	}
	return nil
}

// loadPortfolios reads TOML portfolios config file to cli.portfolios
func (cli *cli) loadPortfolios() error {
	if !fsx.FileExists(cli.configFile) {
		return fmt.Errorf("missing config file: %q", cli.configFile)
	}
	s, err := fsx.ReadFile(cli.configFile)
	if err != nil {
		return err
	}
	conf := struct {
		Portfolios []struct {
			Name        string `toml:"name"`
			Description string `toml:"description"`
			Assets      []struct {
				Symbol      string  `toml:"symbol"`
				Amount      float64 `toml:"amount"`
				Description string  `toml:"description"`
			} `toml:"assets"`
		} `toml:"portfolios"`
	}{}
	_, err = toml.Decode(s, &conf)
	if err != nil {
		return err
	}
	// Copy parsed portfolios configuration to cli.portfolios slice.
	cli.portfolios = portfolio.Portfolios{}
	for _, c := range conf.Portfolios {
		p := portfolio.Portfolio{}
		p.Name = c.Name
		p.Description = c.Description
		p.Assets = []portfolio.Asset{}
		for _, a := range c.Assets {
			asset := portfolio.Asset{}
			asset.Symbol = a.Symbol
			asset.Amount = a.Amount
			asset.Description = a.Description
			p.Assets = append(p.Assets, asset)
		}
		cli.portfolios = append(cli.portfolios, p)
	}
	return err
}
