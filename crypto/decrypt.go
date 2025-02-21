package crypto

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/flashbots/vbackup/aws"
)

func Decrypt(ctx context.Context, key string, data interface{}) (interface{}, error) {
	var dec Decryptor

	switch {
	case aws.IsKMS(key):
		kms, err := aws.NewKMS(key)
		if err != nil {
			return nil, err
		}
		dec = kms
	}

	return decrypt(ctx, dec, data)
}

func decrypt(ctx context.Context, dec Decryptor, data interface{}) (interface{}, error) {
	switch d := data.(type) {

	case map[string]interface{}:
		res := make(map[string]interface{}, len(d))
		for k, v := range d {
			c, err := decrypt(ctx, dec, v)
			if err != nil {
				return nil, err
			}
			if c != nil {
				res[k] = c
			}
		}
		return res, nil

	case []interface{}:
		res := make([]interface{}, len(d))
		for i, v := range d {
			c, err := decrypt(ctx, dec, v)
			if err != nil {
				return nil, err
			}
			if c != nil {
				res[i] = c
			}
		}
		return res, nil

	case string:
		return dec.DecryptString(ctx, d)

	default:
		if data == nil {
			return nil, nil
		}
		fmt.Fprintln(os.Stderr, "WARNING: unexpected type: "+reflect.TypeOf(data).String())
		return nil, nil
	}
}
