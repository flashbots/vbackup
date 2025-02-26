VERSION := $(shell git describe --tags --always --dirty="-dev" --match "v*.*.*" || echo "development" )
VERSION := $(VERSION:v%=%)

.PHONY: build
build:
	@CGO_ENABLED=0 go build \
			-ldflags "-X main.version=${VERSION}" \
			-o ./bin/vbackup \
		github.com/flashbots/vbackup/cmd

.PHONY: snapshot
snapshot:
	@goreleaser release --snapshot --clean

.PHONY: help
help:
	@go run github.com/flashbots/vbackup/cmd/vbackup export --help
	@go run github.com/flashbots/vbackup/cmd/vbackup import --help
