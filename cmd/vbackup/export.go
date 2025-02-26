package main

import (
	"context"
	"errors"
	"os"
	"slices"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/flashbots/vbackup/config"
	"github.com/flashbots/vbackup/crypto"
	"github.com/flashbots/vbackup/vault"
)

func CommandExport(cfg *config.Config) *cli.Command {
	ignore := &cli.StringSlice{}

	flags := slices.Concat(
		keyFlags(cfg),
		vaultFlags(cfg, ignore),
	)

	return &cli.Command{
		Name:  "export",
		Usage: "export vault kv2 store into a file",
		Flags: flags,

		ArgsUsage: " [path/to/local/backup.yaml]",

		Before: func(clictx *cli.Context) error {
			cfg.Vault.Ignore = ignore.Value()

			if clictx.Args().Len() > 1 {
				return errors.New("too many arguments")
			}

			return cfg.Validate()
		},

		Action: func(clictx *cli.Context) error {
			v, err := vault.New(cfg.Vault)
			if err != nil {
				return err
			}
			data, err := v.Get(context.Background(), cfg.Vault.Path)
			if err != nil {
				return err
			}

			if cfg.Key != "" {
				data, err = crypto.Encrypt(context.Background(), cfg.Key, data)
				if err != nil {
					return err
				}
			}

			b, err := yaml.Marshal(data)
			if err != nil {
				return err
			}

			output := os.Stdout
			if clictx.Args().Len() == 1 {
				output, err = os.OpenFile(clictx.Args().First(), os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					return err
				}
				if err := output.Truncate(0); err != nil {
					return err
				}
				defer output.Close()
			}

			if _, err := output.Write(b); err != nil {
				return err
			}

			return nil
		},
	}
}
