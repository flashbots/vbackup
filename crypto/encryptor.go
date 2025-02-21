package crypto

import "context"

type Encryptor interface {
	EncryptString(context.Context, string) (string, error)
}
