package main

import (
	"github.com/flashbots/vbackup/config"
	"github.com/urfave/cli/v2"
)

func keyFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category:    "KEY",
			Destination: &cfg.Key,
			Name:        "key",
			Usage:       "encryption `key`",
		},
	}
}
