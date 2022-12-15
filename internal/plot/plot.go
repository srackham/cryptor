// Plot portfolio data.
package plot

import "github.com/srackham/cryptor/internal/portfolio"

type IPlotter interface {
	PlotValuations(portfolios portfolio.Portfolios) error
	PlotAllocation(portfolio portfolio.Portfolio) error
}

type Plotter struct {
	plotter IPlotter
}

func New(plotter IPlotter) Plotter {
	return Plotter{
		plotter: plotter,
	}
}

func (p *Plotter) PlotValuations(portfolios portfolio.Portfolios) error {
	return p.PlotValuations(portfolios)
}

func (p *Plotter) PlotAllocation(portfolio portfolio.Portfolio) error {
	return p.PlotAllocation(portfolio)
}
