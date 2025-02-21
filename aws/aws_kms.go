package aws

import (
	"context"
	"fmt"
	"os"
	"strings"

	"encoding/base64"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
)

type KMS struct {
	cli *awskms.Client
	key string
}

func IsKMS(id string) bool {
	return strings.HasPrefix(id, "arn:aws:kms:")
}

func NewKMS(id string) (*KMS, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &KMS{
		cli: awskms.NewFromConfig(cfg),
		key: id,
	}, nil
}

func (kms *KMS) EncryptString(ctx context.Context, s string) (string, error) {
	if len(s) > 128*1024 {
		fmt.Fprintf(os.Stderr, "WARNING: strings over 128K are not awskms-encryptable\n")
		return s, nil
	}
	res, err := kms.cli.Encrypt(ctx, &awskms.EncryptInput{
		Plaintext: []byte(s),
		KeyId:     &kms.key,
	})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(res.CiphertextBlob), nil
}

func (kms *KMS) DecryptString(ctx context.Context, s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	res, err := kms.cli.Decrypt(ctx, &awskms.DecryptInput{
		CiphertextBlob: b,
		KeyId:          &kms.key,
	})
	if err != nil {
		return "", err
	}

	return string(res.Plaintext), nil
}
