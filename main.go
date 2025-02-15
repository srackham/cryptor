package main

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/srackham/cryptor/internal/cli"
	. "github.com/srackham/cryptor/internal/global"
	"github.com/srackham/cryptor/internal/helpers"
)

func main() {
	ctx := Context{
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		DataDir:   path.Join(helpers.GetXDGDataDir(), "cryptor"),
		CacheDir:  path.Join(helpers.GetXDGCacheDir(), "cryptor"),
		ConfigDir: path.Join(helpers.GetXDGConfigDir(), "cryptor"),
		Now:       func() time.Time { return time.Now() },
		HttpGet:   http.Get,
	}
	cli := cli.New(&ctx)
	if err := cli.Execute(os.Args...); err != nil {
		os.Exit(1)
	}
}
