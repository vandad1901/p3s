#!/bin/zsh
set -eo pipefail

sudo apt update
sudo apt install -y build-essential wget unzip postgresql-client ca-certificates curl just
