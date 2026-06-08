#!/bin/zsh
set -eo pipefail

VERSION="v1.70.0"

go install github.com/bufbuild/buf/cmd/buf@$VERSION