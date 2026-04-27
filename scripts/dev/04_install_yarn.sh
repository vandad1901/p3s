#!/bin/zsh

set -euo pipefail

export CI=1
export COREPACK_ENABLE_DOWNLOAD_PROMPT=0

npm install -g corepack
corepack enable
corepack prepare yarn@stable --activate
npm install -g @angular/cli --no-fund --no-audit