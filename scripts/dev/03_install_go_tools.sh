#!/bin/zsh
set -eo pipefail

MIGRATE_VERSION="v4.19.1"
PROTOC_GEN_GO_VERSION="v1.36.11"
BUF_VERSION="v1.70.0"

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$MIGRATE_VERSION
go install google.golang.org/protobuf/cmd/protoc-gen-go@$PROTOC_GEN_GO_VERSION
go install github.com/bufbuild/buf/cmd/buf@$BUF_VERSION
