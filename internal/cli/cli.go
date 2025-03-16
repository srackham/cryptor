package cli

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/srackham/cryptor/internal/binance"
	"github.com/srackham/cryptor/internal/fsx"
	. "github.com/srackham/cryptor/internal/global"
	"github.com/srackham/cryptor/internal/portfolio"
	"github.com/srackham/cryptor/internal/slice"
	"github.com/srackham/cryptor/internal/xrates"
	"gopkg.in/yaml.v3"
)

type cli struct {
	*Context
	command     string                // CLI command
	portfolios  portfolio.Portfolios  // Crypto currency portfolios loaded from configuration file
	valuation   portfolio.Portfolios  // Valuated portfolios
	aggregate   portfolio.Portfolio   // Combinded portfolios valuation
	priceReader *binance.PriceReader  // Crypto currency price oracle
	xrates      *xrates.ExchangeRates // Fiat currency to USD exchange rate oracle
	opts        struct {
		aggregate     bool                // Inlcude aggregate (combined) portfolios valuation
		aggregateOnly bool                // Only include aggregate portfolio valuation
		currency      string              // Fiat currency symbol that the valuation is denominated in
		notes         bool                // Include portfolio notes in the valuations
		format        string              // Valuate command output format ("json" or "yaml")
		save          bool                // Update the valuations file
		portfolios    slice.Slice[string] // Names of portfolios to be printed
		prices        portfolio.Prices    // Maps asset symbols to prices
	}
}

// New creates a new cli.
func New(ctx *Context) *cli {
	cli := cli{}
	cli.Context = ctx
	priceReader := binance.NewPriceReader(ctx)
	cli.priceReader = &priceReader
	xrates := xrates.New(ctx)
	cli.xrates = &xrates
	return &cli
}

// Execute runs a command specified by CLI args.
func (cli *cli) Execute(args ...string) error {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintf(cli.Stderr, "\nERROR: %s\n", err.Error())
		}
	}()
	cli.opts.currency = "USD"
	cli.opts.portfolios = slice.New[string]()
	err = cli.parseArgs(args)
	if err != nil {
		return err
	}
	cli.valuation = portfolio.Portfolios{}
	if err != nil {
		return err
	}
	switch cli.command {
	case "help":
		cli.helpCmd()
	case "history":
		err = cli.historyCmd()
	case "init":
		err = cli.initCmd()
	case "valuate":
		err = cli.valuateCmd()
	default:
		err = fmt.Errorf("invalid command: %s", cli.command)
	}
	return err
}

// parseArgs parses and validate command-line arguments.
func (cli *cli) parseArgs(args []string) error {
	skip := false
	cli.opts.prices = make(portfolio.Prices)
	for i, opt := range args {
		if skip {
			skip = false
			continue
		}
		switch {
		case i == 0:
			if len(args) == 1 {
				cli.command = "help"
			}
		case i == 1:
			if opt == "-h" || opt == "--help" {
				opt = "help"
			}
			if !isCommand(opt) {
				return fmt.Errorf("invalid command: \"%s\"", opt)
			}
			cli.command = opt
		case opt == "-aggregate":
			cli.opts.aggregate = true
		case opt == "-aggregate-only":
			cli.opts.aggregateOnly = true
		case opt == "-notes":
			cli.opts.notes = true
		case opt == "-save":
			cli.opts.save = true
		case slice.New("-confdir", "-currency", "-format", "-portfolio", "-price").Has(opt):
			// Process option argument.
			if i+1 >= len(args) {
				return fmt.Errorf("missing %s argument value", opt)
			}
			arg := args[i+1]
			switch opt {
			case "-confdir":
				cli.ConfigDir = arg
				cli.CacheDir = arg
				cli.DataDir = arg
			case "-currency":
				cli.opts.currency = strings.ToUpper(arg)
			case "-format":
				if !slice.New("json", "yaml").Has(arg) {
					return fmt.Errorf("invalid -format argument: \"%s\"", arg)
				}
				cli.opts.format = arg
			case "-portfolio":
				if !portfolio.IsValidName(arg) {
					return fmt.Errorf("invalid -portfolio argument: \"%s\"", arg)
				}
				if cli.opts.portfolios.Has(arg) {
					return fmt.Errorf("-portfolio name can only be specified once: \"%s\"", arg)
				}
				cli.opts.portfolios = append(cli.opts.portfolios, arg)
			case "-price":
				symbol, price, err := ParsePriceOption(arg)
				if err != nil {
					return err
				}
				cli.opts.prices[symbol] = price
			default:
				return fmt.Errorf("unexpected option: \"%s\"", opt)
			}
			skip = true
		default:
			return fmt.Errorf("invalid argument: \"%s\"", opt)
		}
	}
	return nil
}

// ParsePriceOption parses a price option string in the format "SYMBOL=PRICE".
// It returns the uppercase symbol and price as separate values.
// It returns an error if the price option string is invalid or if the symbol or price are invalid.
func ParsePriceOption(priceOption string) (symbol string, price float64, err error) {
	parts := strings.Split(priceOption, "=")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid price option: \"%s\"", priceOption)
	}
	symbol = strings.TrimSpace(parts[0])
	priceStr := strings.TrimSpace(parts[1])
	if !regexp.MustCompile("^[a-zA-Z0-9-_]+$").MatchString(symbol) {
		return "", 0, fmt.Errorf("invalid price symbol: \"%s\"", priceOption)
	}
	price, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid price value: \"%s\"", priceOption)
	}
	return strings.ToUpper(symbol), price, nil
}

// initCmd implements the initCmd command.
func (cli *cli) initCmd() error {
	if !fsx.DirExists(cli.ConfigDir) {
		fmt.Fprintf(cli.Stdout, "creating configuration directory: \"%s\"\n", cli.ConfigDir)
		if err := fsx.MkMissingDir(cli.ConfigDir); err != nil {
			return err
		}
	}
	if fsx.FileExists(cli.configFile()) {
		fmt.Fprintf(cli.Stdout, "config file already exists: \"%s\"\n", cli.configFile())
	} else {
		fmt.Fprintf(cli.Stdout, "installing example config file: \"%s\"\n", cli.configFile())
		contents := `# NOTE: The xrates-appid option is only necessary if the -currency command option is used.
# Open Exchange Rates App ID (https://openexchangerates.org/)
xrates-appid: YOUR_APP_ID`
		if err := fsx.WriteFile(cli.configFile(), contents); err != nil {
			return fmt.Errorf("failed to write config file: \"%s\"", err.Error())
		}
	}
	if fsx.FileExists(cli.portfoliosFile()) {
		fmt.Fprintf(cli.Stdout, "portfolios file already exists: \"%s\"\n", cli.portfoliosFile())
	} else {
		fmt.Fprintf(cli.Stdout, "installing example portfolios file: \"%s\"\n", cli.portfoliosFile())
		contents := `# Example cryptor portfolios configuration file containing two portfolios installed by 'cryptor init' command.

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
    BTC: 1.0`
		if err := fsx.WriteFile(cli.portfoliosFile(), contents); err != nil {
			return err
		}
	}
	if !fsx.DirExists(cli.CacheDir) {
		fmt.Fprintf(cli.Stdout, "creating cache directory: \"%s\"\n", cli.CacheDir)
		if err := fsx.MkMissingDir(cli.CacheDir); err != nil {
			return err
		}
	}
	if !fsx.DirExists(cli.DataDir) {
		fmt.Fprintf(cli.Stdout, "creating data directory: \"%s\"\n", cli.DataDir)
		if err := fsx.MkMissingDir(cli.DataDir); err != nil {
			return err
		}
	}
	return nil
}

// historyCmd prints the saved valuations history.
func (cli *cli) historyCmd() (err error) {
	fname := cli.valuationsFile("json")
	valuations := portfolio.Portfolios{}
	if fsx.FileExists(fname) {
		valuations, err = portfolio.LoadValuations(fname)
		if err != nil {
			return fmt.Errorf("valuations file: \"%s\": %s", fname, err.Error())
		}
	}
	if len(cli.opts.portfolios) > 0 {
		valuations = valuations.FilterByName(cli.opts.portfolios...)
	}
	if len(valuations) == 0 {
		return fmt.Errorf("valuations file: \"%s\": no valuations found", fname)
	}
	var s string
	if cli.opts.format == "yaml" {
		s, err = valuations.ToYAML()
	} else {
		s, err = valuations.ToJSON()
	}
	if err == nil {
		_, err = fmt.Fprint(cli.Stdout, s)
	}
	return
}

// helpCmd implements the `help` command.
func (cli *cli) helpCmd() {
	github := "https://github.com/srackham/cryptor"
	summary := `
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
    -portfolio PORTFOLIO        Process named portfolio (default: all portfolios)
    -price SYMBOL=PRICE         Override the asset price of SYMBOL with PRICE (in USD)
    -format FORMAT              Set the valuate command output format ("json" or "yaml")

Config directory: ` + cli.ConfigDir + `
Cache directory:  ` + cli.CacheDir + `
Data directory:   ` + cli.DataDir + `

Version:    ` + VERS + " (" + OS + ")" + `
Git commit: ` + COMMIT + `
Built:      ` + BUILT + `
Github:     ` + github

	fmt.Fprintf(cli.Stdout, "%s\n", summary)
}

func isCommand(name string) bool {
	return slice.New("help", "history", "init", "valuate").Has(name)
}

func (cli *cli) configFile() string {
	return filepath.Join(cli.ConfigDir, "config.yaml")
}

func (cli *cli) portfoliosFile() string {
	return filepath.Join(cli.ConfigDir, "portfolios.yaml")
}

func (cli *cli) valuationsFile(format string) string {
	return filepath.Join(cli.DataDir, "valuations."+format)
}

func (cli *cli) loadPortfolios() (err error) {
	ps, err := cli.loadConfigFile(cli.portfoliosFile())
	if err != nil {
		return fmt.Errorf("portfolios file: \"%s\": %s", cli.portfoliosFile(), err.Error())
	}
	if err = ps.Validate(true); err != nil {
		return fmt.Errorf("portfolios file: \"%s\": %s", cli.portfoliosFile(), err.Error())
	}
	cli.portfolios = ps
	return nil
}

// save appends the current valuation to the valuations file and saves the exchange rates cache file.
func (cli *cli) save() (err error) {
	if cli.opts.save {
		fname := cli.valuationsFile("json")
		valuations := portfolio.Portfolios{}
		if fsx.FileExists(fname) {
			valuations, err = portfolio.LoadValuations(fname)
			if err != nil {
				return fmt.Errorf("valuations file: \"%s\": %s", fname, err.Error())
			}
		}
		valuations = append(valuations, cli.valuation...)
		valuations = append(valuations, cli.aggregate)
		err = valuations.SaveValuations(fname)
		if err != nil {
			return fmt.Errorf("valuations file: \"%s\": %s", fname, err.Error())
		}
	}
	if len(*(cli.xrates.CacheData)) > 0 {
		err = cli.xrates.Save(cli.xrates.CacheFile())
		if err != nil {
			return fmt.Errorf("exchange rates file: \"%s\": %s", cli.xrates.CacheFile(), err.Error())
		}
	}
	return
}

// valuateCmd implements the valuate command.
func (cli *cli) valuateCmd() error {
	now := cli.Now()
	date := now.Format("2006-01-02")
	time := now.Format("15:04:05")
	if err := cli.loadPortfolios(); err != nil {
		return err
	}
	// Select portfolios to be valuated.
	cli.valuation = portfolio.Portfolios{}
	cli.valuation = cli.portfolios
	// Evaluate portfolios.
	for i := range cli.valuation {
		cli.valuation[i].Date = date
		cli.valuation[i].Time = time
		if err := cli.valuation[i].SetUSDValues(cli.priceReader); err != nil {
			return err
		}
		cli.valuation[i].SetAllocations()
		cli.valuation[i].Assets.Sort()
	}
	cli.aggregate = cli.valuation.Aggregate("aggregate")
	cli.aggregate.Date = date
	cli.aggregate.Time = time
	printed_valuation := cli.valuation
	if len(cli.opts.portfolios) > 0 {
		// Select -portfolio option valuations.
		printed_valuation = portfolio.Portfolios{}
		for _, p := range cli.valuation {
			if cli.opts.portfolios.IndexOf(p.Name) >= 0 {
				printed_valuation = append(printed_valuation, p)
			}
		}
	}
	if cli.opts.aggregateOnly {
		printed_valuation = portfolio.Portfolios{cli.aggregate}
	} else if cli.opts.aggregate {
		printed_valuation = append(printed_valuation, cli.aggregate)
	}
	// Print portfolios.
	xrate, err := cli.xrates.GetCachedRate(cli.opts.currency, false)
	if err != nil {
		return err
	}
	if s, err := printed_valuation.ToString(cli.opts.format, cli.opts.currency, xrate); err != nil {
		return err
	} else {
		fmt.Fprintf(cli.Stdout, "\n%s\n", s)
	}
	// Save valuation and cache files.
	if err := cli.save(); err != nil {
		return err
	}
	return nil
}

// loadConfigFile reads portfolios configuration file.
func (cli *cli) loadConfigFile(filename string) (portfolio.Portfolios, error) {
	type Config []struct {
		Name   string             `yaml:"name"`
		Notes  string             `yaml:"notes"`
		Cost   string             `yaml:"cost"`
		Assets map[string]float64 `yaml:"assets"`
	}
	res := portfolio.Portfolios{}
	s, err := fsx.ReadFile(filename)
	if err != nil {
		return res, err
	}
	config := Config{}
	err = yaml.Unmarshal([]byte(s), &config)
	if err != nil {
		// Try to parse a minimal portfolio
		assets := make(map[string]float64)
		err = yaml.Unmarshal([]byte(s), &assets)
		if err != nil {
			return res, err
		}
		config = Config{{
			Name:   "portfolio1",
			Assets: assets,
		}}
	}
	// Copy parsed portfolios configuration to Portfolios slice.
	for _, c := range config {
		p := portfolio.Portfolio{}
		p.Name = c.Name
		if cli.opts.notes {
			p.Notes = c.Notes
		}
		p.Assets = []portfolio.Asset{}
		for k, v := range c.Assets {
			asset := portfolio.Asset{}
			asset.Symbol = strings.ToUpper(k)
			asset.Amount = v
			p.Assets = append(p.Assets, asset)
		}
		res = append(res, p)
	}
	// Synthesise missing portfolio names.
	for i := range res {
		if res[i].Name == "" {
			for j := 1; ; j++ {
				name := fmt.Sprintf("portfolio%d", j)
				if res.FindByName(name) == -1 {
					res[i].Name = name
					break
				}
			}
		}
	}
	// Check for duplicate portfolio names.
	for i := range res {
		for j := range res {
			if i != j && res[i].Name == res[j].Name {
				return res, fmt.Errorf("duplicate portfolio name: \"%s\"", res[j].Name)
			}
		}
	}
	// Assign asset price options
	for symbol, price := range cli.opts.prices {
		if err := res.SetAssetPrice(symbol, price); err != nil {
			return res, err
		}
	}
	// Check all -portfolio option names exist.
	for _, name := range cli.opts.portfolios {
		i := res.FindByName(name)
		if i == -1 {
			return res, fmt.Errorf("missing portfolio: \"%s\"", name)
		}
	}
	// Calculate portfolio costs in USD (this is done last to avoid loading the exchange rates cache unnecessarily)
	for i := range config {
		if config[i].Cost != "" {
			if config[i].Name != res[i].Name {
				panic("out of order portfolios")
			}
			usd, err := cli.currencyToUSD(config[i].Cost)
			if err != nil {
				return res, err
			}
			res[i].Cost = usd
		}
	}
	return res, err
}

// currencyToUSD parses "<value>[<currency]" string and converts to USD.
func (cli *cli) currencyToUSD(currencyValue string) (value float64, err error) {
	value, currency, err := portfolio.ParseCurrency(currencyValue)
	if err != nil {
		return
	}
	rate, err := cli.xrates.GetCachedRate(currency, false)
	if err != nil {
		return
	}
	if currency == "USD" && rate != 1.00 {
		err = fmt.Errorf("USD exchange rate should be 1.00: %f", rate)
		return
	}
	if rate <= 0.00 {
		err = fmt.Errorf("exchange rate is zero or less: %.2f %s", rate, currency)
		return
	}
	value = value / rate
	return
}
