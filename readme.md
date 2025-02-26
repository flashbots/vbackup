# vbackup

Import/export (encrypted) values between vault kv2 path and a YAML file.

## TL;DR

```shell
AWS_PROFILE=xxxx \
go run github.com/flashbots/vbackup/cmd/vbackup export \
    --key arn:aws:kms:us-east-2:xxxxxxxxxxxx:key/yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy \
    --mount secret/kv \
    --path path/to/some/nested/secret \
    --ignore subpath/you/want/to/ignore,another/subpath/you/want/to/ignore \
  vault.yaml
```

```shell
AWS_PROFILE=xxxx \
go run github.com/flashbots/vbackup/cmd/vbackup import \
    --dry-run \
    --key arn:aws:kms:us-east-2:xxxxxxxxxxxx:key/yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy \
    --mount secret/kv \
    --path path/to/some/nested/secret \
  vault.yaml
```
