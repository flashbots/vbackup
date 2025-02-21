package crypto

import "context"

type Decryptor interface {
	DecryptString(context.Context, string) (string, error)
}
