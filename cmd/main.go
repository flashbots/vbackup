package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/flashbots/vbackup/config"
)

var (
	version = "development"
)

func main() {
	cfg := config.New()

	flags := []cli.Flag{}

	commands := []*cli.Command{
		CommandExport(cfg),
		CommandImport(cfg),
		CommandHelp(cfg),
	}

	app := &cli.App{
		Name:    "vbackup",
		Usage:   "Export/import kv store contents from hashicorp vault",
		Version: version,

		Flags:    flags,
		Commands: commands,

		Before: func(_ *cli.Context) error {
			return nil
		},

		Action: func(clictx *cli.Context) error {
			return cli.ShowAppHelp(clictx)
		},
	}

	defer func() {
		zap.L().Sync() //nolint:errcheck
	}()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\nFailed with error:\n\n%s\n\n", err.Error())
		os.Exit(1)
	}
}
