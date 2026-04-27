#!/bin/zsh
set -eo pipefail

VERSION="v4.19.1"

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$VERSION