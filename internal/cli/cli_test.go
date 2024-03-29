package cli

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/fsx"
	"github.com/srackham/cryptor/internal/helpers"
	"github.com/srackham/cryptor/internal/mockprice"
)

func mockCli() *cli {
	return New(&mockprice.Reader{})
}

func TestParseArgs(t *testing.T) {
	var c cli
	var err error
	parse := func(cmd string) {
		args := strings.Split(cmd, " ")
		c = *mockCli()
		err = c.parseArgs(args)
	}

	parse("cryptor")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, "help", c.command)
	parse("cryptor invalid-command")
	assert.Equal(t, `invalid command: "invalid-command"`, err.Error())
}

func exec(cli *cli, cmd string) (out string, err error) {
	tmpdir, err := os.MkdirTemp("", "cryptor")
	if err != nil {
		return
	}
	if !helpers.GithubActions() {
		err = fsx.CopyFile("../../testdata/config.yaml", path.Join(tmpdir, "config.yaml"))
		if err != nil {
			return
		}
	}
	err = fsx.CopyFile("../../testdata/portfolios.yaml", path.Join(tmpdir, "portfolios.yaml"))
	if err != nil {
		return
	}
	cmd = fmt.Sprintf("%s -confdir %s", cmd, tmpdir)
	args := strings.Split(cmd, " ")
	cli.log.Out = make(chan string, 100)
	err = cli.Execute(args)
	close(cli.log.Out)
	if err != nil {
		return
	}
	for line := range cli.log.Out {
		out += line + "\n"
	}
	out = strings.Replace(out, `\`, `/`, -1) // Normalize MS Windows path separators.
	return
}

func TestEvaluateCmd(t *testing.T) {
	today := helpers.TodaysDate()
	cli := mockCli()
	out, err := exec(cli, "cryptor valuate")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "NAME:  personal\nNOTES: ## Personal Portfolio\n")
	fmt.Println(out)
	assert.PassIf(t, cli.valuations.FindByNameAndDate("personal", today) != -1, "missing valuation: %v", today)

	cli = mockCli()
	out, err = exec(cli, "cryptor valuate -date "+today+" -force")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "price request: BTC "+today+" 10000.00")
	assert.Contains(t, out, "price request: ETH "+today+" 1000.00")
	assert.Contains(t, out, "price request: USDC "+today+" 1.00")
	fmt.Println(out)

	cli = mockCli()
	out, err = exec(cli, "cryptor valuate -date 0 -force")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "price request: BTC "+today+" 10000.00")
	assert.Contains(t, out, "price request: ETH "+today+" 1000.00")
	assert.Contains(t, out, "price request: USDC "+today+" 1.00")
	fmt.Println(out)

	cli = mockCli()
	date := "2022-06-01"
	out, err = exec(cli, "cryptor valuate -date "+date)
	assert.Contains(t, out, "price request: BTC "+date+" 10000.00")
	assert.Contains(t, out, "price request: ETH "+date+" 1000.00")
	assert.Contains(t, out, "price request: USDC "+date+" 1.00")
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, cli.valuations.FindByNameAndDate("personal", date) != -1, "missing personal valuation: %v", date)
	assert.PassIf(t, cli.valuations.FindByNameAndDate("joint", date) != -1, "missing joint valuation: %v", date)
	assert.PassIf(t, cli.valuations.FindByNameAndDate("portfolio1", date) != -1, "missing portfolio1 valuation: %v", date)
}

func TestHelpCmd(t *testing.T) {
	out, err := exec(mockCli(), "cryptor help")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "Cryptor valuates")

	out, err = exec(mockCli(), "cryptor --help")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "Cryptor valuates")

	out, err = exec(mockCli(), "cryptor -h")
	assert.PassIf(t, err == nil, "%v", err)
	assert.Contains(t, out, "Cryptor valuates")

}
