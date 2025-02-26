package main

import (
	"context"
	"errors"
	"io"
	"os"
	"slices"

	"github.com/flashbots/vbackup/config"
	"github.com/flashbots/vbackup/crypto"
	"github.com/flashbots/vbackup/vault"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func CommandImport(cfg *config.Config) *cli.Command {
	dryRun := false
	ignore := &cli.StringSlice{}

	importFlags := []cli.Flag{
		&cli.BoolFlag{
			Destination: &dryRun,
			Name:        "dry-run",
			Usage:       "dry run (decrypt and match, but don't put/patch vault)",
			Value:       false,
		},
	}

	flags := slices.Concat(
		importFlags,
		keyFlags(cfg),
		vaultFlags(cfg, ignore),
	)

	return &cli.Command{
		Name:  "import",
		Usage: "import vault kv2 store from a file",
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
			input := os.Stdin
			if clictx.Args().Len() == 1 {
				_input, err := os.Open(clictx.Args().First())
				if err != nil {
					return err
				}
				input = _input
				defer _input.Close()
			}

			b, err := io.ReadAll(input)
			if err != nil {
				return err
			}

			data := make(map[string]interface{})
			if err := yaml.Unmarshal(b, &data); err != nil {
				return err
			}

			if cfg.Key != "" {
				_data, err := crypto.Decrypt(context.Background(), cfg.Key, data)
				if err != nil {
					return err
				}
				var ok bool
				data, ok = _data.(map[string]interface{})
				if !ok {
					panic("decrypt returned wrong value type")
				}
			}

			v, err := vault.New(cfg.Vault)
			if err != nil {
				return err
			}
			if err := v.Put(context.Background(), cfg.Vault.Path, data, dryRun); err != nil {
				return err
			}

			return nil
		},
	}
}
