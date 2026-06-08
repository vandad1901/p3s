#!/bin/zsh
set -eo pipefail

VERSION="v1.36.11"

go install google.golang.org/protobuf/cmd/protoc-gen-go@$VERSION