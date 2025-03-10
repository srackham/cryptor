package portfolio

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/srackham/cryptor/internal/binance"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/set"
	"gopkg.in/yaml.v3"
)

// An amount of crypto currency belonging to a portfolio.
type Asset struct {
	Symbol     string  `yaml:"symbol"     json:"symbol"`     // Crypto currecy symbol
	Amount     float64 `yaml:"amount"     json:"amount"`     // Numerical amount crypto currency
	Value      float64 `yaml:"value"      json:"value"`      // Asset value in USD at the time of valuation
	Allocation float64 `yaml:"allocation" json:"allocation"` // Percentage of total portfolio value
}

type Assets []Asset

// Portfolio stores a named portfolio of zero or more crypto currency assets.
// - A Portfolio is loaded from a `portfolios.yaml` configuration file.
// - The `Valuate` method calculates and sets the current value of a `Portfolio`.
// - The valuated `Portfolio` is appended to a `valuations.yaml` file.
// - Note that the portfolios configuration and valuations files have different formats.
type Portfolio struct {
	Name    string  `yaml:"name"     json:"name"`  // Porfolio name
	Notes   string  `yaml:"notes"    json:"notes"` // User notes
	Date    string  `yaml:"date"     json:"date"`  // The valuation date formatted "YYYY-MM-DD"
	Time    string  `yaml:"time"     json:"time"`  // The valuation time formatted "hh:mm:ss""
	Value   float64 `yaml:"value"    json:"value"` // Current portfolio value in USD
	Paid    string  `yaml:"paid"     json:"paid"`  // The amount paid for the portfolio from portfolios.yaml configuration file
	PaidUSD float64 `yaml:"cost"     json:"cost"`  // The amount paid for the portfolio calculated in USD at the current exchange rate
	Assets  Assets  `yaml:"assets"   json:"assets"`
}

type Portfolios []Portfolio

// Returns `true` if the portfolio `name` is valid.
func IsValidName(name string) bool {
	re := regexp.MustCompile(`^\w[-\w]*$`)
	return re.MatchString(name)
}

// ParseCurrency extracts the amount and currency symbol from a string formatted like "<amount><symbol>".
// <symbol> defaults to "USD".
func ParseCurrency(currencyValue string) (amount float64, symbol string, err error) {
	s := strings.ReplaceAll(currencyValue, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	re := regexp.MustCompile(`^([0-9.]+)([a-zA-Z]*)$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		err = fmt.Errorf("invalid currency value: \"%s\"", currencyValue)
		return
	}
	amountStr := matches[1]
	symbolStr := matches[2]
	amount, err = strconv.ParseFloat(amountStr, 64)
	if err != nil {
		err = fmt.Errorf("invalid currency value: \"%s\"", currencyValue)
		return
	}
	if symbolStr == "" {
		symbol = "USD"
	} else {
		symbol = strings.ToUpper(symbolStr)
	}
	return
}

// Sort assets by descending value.
func (assets Assets) Sort() {
	// TODO tests
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Value > assets[j].Value
	})
}

// Find searches assets slice for an asset with matching `symbol`.
// If found it return the asset index else returns -1.
func (assets Assets) Find(symbol string) int {
	// TODO tests
	for i := range assets {
		if assets[i].Symbol == symbol {
			return i
		}
	}
	return -1
}

// SetUSDValues calculates the current USD value of portfolio assets and their total value.
// TODO tests
func (p *Portfolio) SetUSDValues(reader *binance.PriceReader) error {
	total := 0.0
	for i, a := range p.Assets {
		rate, err := reader.GetCachedPrice(a.Symbol)
		if err != nil {
			return err
		}
		val := a.Amount * rate
		p.Assets[i].Value = val
		total += val
	}
	p.Value = total
	return nil
}

// SetAllocations synthesizes asset allocation as a percentage of the total portfolio value.
// TODO tests
func (p *Portfolio) SetAllocations() {
	for i, a := range p.Assets {
		if p.Value != 0.00 {
			p.Assets[i].Allocation = a.Value / p.Value * 100
		}
	}
}

// See [How to deep copy a struct in Go](https://www.educative.io/answers/how-to-deep-copy-a-struct-in-go)
// TODO tests
func (p Portfolio) DeepCopy() Portfolio {
	res := p
	res.Assets = nil
	res.Assets = append(res.Assets, p.Assets...)
	return res
}

func (ps Portfolios) ToJSON() (string, error) {
	data, err := json.MarshalIndent(ps, "", "  ")
	return string(data) + "\n", err
}

func (ps Portfolios) ToYAML() (string, error) {
	data, err := yaml.Marshal(ps)
	return string(data), err
}

// LoadValuations reads a file of portfolio valuations.
func LoadValuations(fname string) (Portfolios, error) {
	res := Portfolios{}
	format := strings.ToLower(filepath.Ext(fname)[1:])
	s, err := fsx.ReadFile(fname)
	if err == nil {
		switch format {
		case "json":
			err = json.Unmarshal([]byte(s), &res)
		case "yaml":
			err = yaml.Unmarshal([]byte(s), &res)
		default:
			err = fmt.Errorf("invalid format: \"%s\"", format)
		}
	}
	return res, err
}

// SaveValuations appends the valuated portfolios to file `fname` in JSON or YAML `format.
func (ps Portfolios) SaveValuations(fname string) (err error) {
	valuations := Portfolios{}
	if fsx.FileExists(fname) {
		valuations, err = LoadValuations(fname)
		if err != nil {
			return
		}
	}
	valuations = append(valuations, ps...)
	format := strings.ToLower(filepath.Ext(fname)[1:])
	var s string
	switch format {
	case "json":
		s, err = ps.ToJSON()
	case "yaml":
		s, err = ps.ToYAML()
	default:
		err = fmt.Errorf("invalid format: \"%s\"", format)
	}
	if err != nil {
		return
	}
	err = fsx.WriteFile(fname, s)
	return
}

// Aggregate returns a new portfolio that combines the valuated receiver portfolios.
// Portfolio Notes field is assigned the list of combined portfolios.
// Aggregated costs are valid only if all portfolios are costed.
// TODO tests
func (ps Portfolios) Aggregate(name string) Portfolio {
	res := Portfolio{
		Name:   name,
		Assets: Assets{},
	}
	notes := []string{}
	isMissingCost := false
	for _, p := range ps {
		notes = append(notes, p.Name)
		res.Value += p.Value
		if p.PaidUSD == 0 {
			isMissingCost = true
		}
		res.PaidUSD += p.PaidUSD
		for _, a := range p.Assets {
			i := res.Assets.Find(a.Symbol)
			if i == -1 {
				res.Assets = append(res.Assets, Asset{Symbol: a.Symbol, Amount: a.Amount, Value: a.Value})
			} else {
				res.Assets[i].Amount += a.Amount
				res.Assets[i].Value += a.Value
			}
		}
	}
	sort.Strings(notes)
	res.Notes = strings.Join(notes, ", ")
	res.SetAllocations()
	res.Assets.Sort()
	if isMissingCost {
		res.PaidUSD = 0.00 // Cost is "omitted" if one or more portfolios are not costed
	}
	return res
}

// FilterByDate returns a list of portfolios dated `date`.
func (ps Portfolios) FilterByDate(date string) Portfolios {
	res := []Portfolio{}
	for _, p := range ps {
		if p.Date == date {
			res = append(res, p)
		}
	}
	return res
}

// FindByNameAndDate searches portfolios slice for a portfolio whose name and date matches portfolio `p`.
// If found it return the portfolio index else returns -1.
// TODO tests
func (ps Portfolios) FindByNameAndDate(name string, date string) int {
	for i := range ps {
		if name == ps[i].Name && date == ps[i].Date {
			return i
		}
	}
	return -1
}

// Find searches portfolios slice for a portfolio whose name matches `name`.
// If found return the portfolio index else return -1.
// TODO tests
func (ps Portfolios) FindByName(name string) int {
	for i := range ps {
		if name == ps[i].Name {
			return i
		}
	}
	return -1
}

// Validate returns `true` if the portfolio fields pass basic sanity checks.
// If `nodups` is `true` duplicate portfolio names are disallowed.
// TODO tests
func (ps Portfolios) Validate(nodups bool) error {
	names := set.New[string]()
	for _, p := range ps {
		if !IsValidName(p.Name) {
			return fmt.Errorf("invalid portfolio name: \"%s\"", p.Name)
		}
		assets := set.New[string]()
		for _, a := range p.Assets {
			if !IsValidName(a.Symbol) {
				return fmt.Errorf("invalid portfolio asset name: \"%s\"", a.Symbol)
			}
			if assets.Has(a.Symbol) {
				return fmt.Errorf("duplicate asset name: \"%s\"", a.Symbol)
			}
			assets.Add(a.Symbol)
		}
		if nodups {
			if names.Has(p.Name) {
				return fmt.Errorf("duplicate portfolio name: \"%s\"", p.Name)
			}
			names.Add(p.Name)
		}
	}
	return nil
}

// FilterByName returns a list of named portfolios.
func (ps Portfolios) FilterByName(names ...string) Portfolios {
	res := []Portfolio{}
	for _, p := range ps {
		for _, name := range names {
			if p.Name == name {
				res = append(res, p)
			}
		}
	}
	return res
}

// Sort Portfolios by ascending date, time and name.
func (ps Portfolios) Sort() {
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].Date != ps[j].Date {
			return ps[i].Date < ps[j].Date
		}
		if ps[i].Time != ps[j].Time {
			return ps[i].Time < ps[j].Time
		}
		return ps[i].Name < ps[j].Name
	})
}

func (p Portfolio) gains() float64 {
	return p.Value - p.PaidUSD
}

func (p Portfolio) pcgains() float64 {
	if p.PaidUSD > 0.00 {
		return p.gains() / p.PaidUSD * 100
	} else {
		return 0.0
	}
}

func (ps *Portfolios) ToText(currency string, xrate float64) string {
	res := ""
	for _, p := range *ps {
		res += fmt.Sprintf("NAME:  %s\nNOTES: %s\nDATE:  %s\nTIME:  %s\nVALUE: %.2f %s",
			p.Name, p.Notes, p.Date, p.Time, p.Value*xrate, currency)
		if p.PaidUSD > 0.00 {
			res += fmt.Sprintf("\nPAID:  %.2f %s\nGAINS: %.2f %s (%.2f%%)", p.PaidUSD*xrate, currency, p.gains()*xrate, currency, p.pcgains())
		} else {
			res += "\nPAID:\nGAINS:"
		}
		if currency != "USD" {
			res += fmt.Sprintf("\nXRATE: 1 USD = %.2f %s", xrate, currency)
		} else {
			res += "\nXRATE:"
		}
		res += "\n            AMOUNT            VALUE    PERCENT       UNIT PRICE\n"
		for _, a := range p.Assets {
			value := a.Value * xrate
			res += fmt.Sprintf("%-5s %12.4f %12.2f %s    %6.2f%% %12.2f %s\n",
				a.Symbol,
				a.Amount,
				value,
				currency,
				a.Allocation,
				helpers.If(a.Amount > 0.0, value/a.Amount, 0),
				currency)
		}
		res += "\n"
	}
	return res
}

func (ps Portfolios) ToString(format string, currency string, xrate float64) (res string, err error) {
	switch format {
	case "":
		res = ps.ToText(currency, xrate)
	case "json":
		res, err = ps.ToJSON()
		if err != nil {
			return
		}
	case "yaml":
		res, err = ps.ToYAML()
		if err != nil {
			return
		}
	default:
		panic(fmt.Sprintf("invalid format: \"%s\"", format))
	}
	res = strings.TrimSpace(res)
	return
}
