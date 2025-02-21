package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/flashbots/vbackup/aws"
)

func Encrypt(ctx context.Context, key string, data interface{}) (interface{}, error) {
	var enc Encryptor

	switch {
	case aws.IsKMS(key):
		kms, err := aws.NewKMS(key)
		if err != nil {
			return nil, err
		}
		enc = kms
	}

	return encrypt(ctx, enc, data)
}

func encrypt(ctx context.Context, enc Encryptor, data interface{}) (interface{}, error) {
	switch d := data.(type) {

	case map[string]interface{}:
		res := make(map[string]interface{}, len(d))
		for k, v := range d {
			c, err := encrypt(ctx, enc, v)
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
			c, err := encrypt(ctx, enc, v)
			if err != nil {
				return nil, err
			}
			if c != nil {
				res[i] = c
			}
		}
		return res, nil

	case string:
		return enc.EncryptString(ctx, d)

	case int:
		return enc.EncryptString(ctx, strconv.Itoa(d))

	case json.Number:
		return enc.EncryptString(ctx, d.String())

	default:
		if data == nil {
			return nil, nil
		}
		fmt.Fprintln(os.Stderr, "WARNING: unexpected type: "+reflect.TypeOf(data).String())
		return nil, nil
	}
}
