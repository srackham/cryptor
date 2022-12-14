package main

import (
	"os"

	"github.com/srackham/cryptor/internal/cgprice"
	"github.com/srackham/cryptor/internal/cli"
)

func main() {
	cli := cli.New(cgprice.NewReader())
	if err := cli.Execute(os.Args); err != nil {
		os.Exit(1)
	}
}
