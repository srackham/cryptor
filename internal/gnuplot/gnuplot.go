// Plot using gnuplot (http://www.gnuplot.info/).
// Implements interface IPlotter.

package gnuplot

import "github.com/srackham/cryptor/internal/portfolio"

type Plotter struct{}

func (p *Plotter) PlotValuations(portfolios portfolio.Portfolios) error {
	// TODO
	return nil
}

func (p *Plotter) PlotAllocation(portfolio portfolio.Portfolio) error {
	// TODO
	return nil
}
