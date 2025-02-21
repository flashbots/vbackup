package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	vapi "github.com/hashicorp/vault/api"
)

func (v *Vault) Put(
	ctx context.Context,
	path string,
	data interface{},
	dryRyn bool,
) error {
	return v.put(ctx, path, "", data, dryRyn)
}

func (v *Vault) put(
	ctx context.Context,
	path string,
	breadcrumbs string,
	data interface{},
	dryRyn bool,
) error {
	// traverse recursively until we reach "leaf" root (a node that only has atomic children)
	if !isLeafRoot(data) {
		switch d := data.(type) {

		case map[string]interface{}:
			for key, val := range d {
				key = sanitised(key)
				if _, ignore := v.ignore[breadcrumbs+"/"+key]; ignore {
					continue
				}
				if err := v.put(ctx, path+"/"+key, breadcrumbs+"/"+key, val, dryRyn); err != nil {
					return err
				}
			}

		default:
			if data != nil {
				fmt.Fprintf(os.Stderr, "WARNING: unexpected type at %s: %s\n",
					path,
					reflect.TypeOf(data).String(),
				)
			}
		}

		return nil
	}

	// write out the leaf root
	switch d := data.(type) {

	case map[string]interface{}:
		s, err := v.kvv2.Get(ctx, path)
		if err != nil && !errors.Is(err, vapi.ErrSecretNotFound) {
			return err
		}

		if s == nil || s.Data == nil {
			fmt.Fprintf(os.Stderr, "writing: %s\n", path)
			if !dryRyn {
				if _, err := v.kvv2.Put(ctx, path, d); err != nil {
					return err
				}
			}
		} else {
			p := make(map[string]interface{}, len(d))
			for key, newValue := range d {
				oldValue, exists := s.Data[key]
				if !exists || !equal(newValue, oldValue) {
					p[key] = newValue
				}
			}
			if len(p) > 0 {
				fmt.Fprintf(os.Stderr, "patching: %s\n", path)
				if !dryRyn {
					if _, err := v.kvv2.Patch(ctx, path, p); err != nil {
						return err
					}
				}
			} else {
				fmt.Fprintf(os.Stderr, "skipping as already matching: %s\n", path)
			}
		}

	default:
		if data != nil {
			fmt.Fprintf(os.Stderr, "WARNING: unexpected type at %s: %s\n",
				path,
				reflect.TypeOf(data).String(),
			)
		}
	}

	return nil
}
