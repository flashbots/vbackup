package main

import (
	"github.com/flashbots/vbackup/config"
	"github.com/urfave/cli/v2"
)

func vaultFlags(cfg *config.Config, ignore *cli.StringSlice) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category:    "VAULT",
			Destination: &cfg.Vault.Address,
			Name:        "addr",
			Usage:       "`url` address of vault",
		},

		&cli.StringFlag{
			Category:    "VAULT",
			Destination: &cfg.Vault.Mount,
			Name:        "mount",
			Usage:       "`path` where kv engine is mounted",
			Value:       "kv",
		},

		&cli.StringFlag{
			Category:    "VAULT",
			Destination: &cfg.Vault.Path,
			Name:        "path",
			Usage:       "`path` to the secrets",
			Value:       "/",
		},

		&cli.StringSliceFlag{
			Category:    "VAULT",
			Destination: ignore,
			Name:        "ignore",
			Usage:       "list of `subpaths` to ignore",
		},
	}
}
