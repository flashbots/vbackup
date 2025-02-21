package vault

import (
	"github.com/flashbots/vbackup/config"
	vapi "github.com/hashicorp/vault/api"
	vconfig "github.com/hashicorp/vault/api/cliconfig"
)

type Vault struct {
	cli     *vapi.Client
	kvv2    *vapi.KVv2
	logical *vapi.Logical

	mount string

	ignore map[string]struct{}
}

func New(cfg *config.Vault) (*Vault, error) {
	config := vapi.DefaultConfig()

	if cfg.Address != "" {
		config.Address = cfg.Address
	}

	cli, err := vapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	if cli.Token() == "" {
		helper, err := vconfig.DefaultTokenHelper()
		if err != nil {
			return nil, err
		}
		token, err := helper.Get()
		if err != nil {
			return nil, err
		}
		cli.SetToken(token)
	}

	v := &Vault{
		cli:     cli,
		kvv2:    cli.KVv2(cfg.Mount),
		logical: cli.Logical(),
		mount:   cfg.Mount,
	}

	if len(cfg.Ignore) > 0 {
		v.ignore = make(map[string]struct{}, len(cfg.Ignore))
		for _, subpath := range cfg.Ignore {
			v.ignore[subpath] = struct{}{}
		}
	}

	return v, nil
}
