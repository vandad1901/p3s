#!/bin/zsh
set -eo pipefail

ASDF_VERSION="0.18.1"
TMP_DIR="$(mktemp -d)"
ZSHRC="$HOME/.zshrc"

wget -qO- "https://github.com/asdf-vm/asdf/releases/download/v$ASDF_VERSION/asdf-v$ASDF_VERSION-linux-amd64.tar.gz" | tar -xz -C "$TMP_DIR"

sudo mv "$TMP_DIR/asdf" /usr/local/bin/asdf
sudo chmod +x /usr/local/bin/asdf
rm -rf "$TMP_DIR"

SHIMS_LINE='export PATH="${ASDF_DATA_DIR:-$HOME/.asdf}/shims:$PATH"'
if ! grep -Fxq "$SHIMS_LINE" "$ZSHRC"; then
  echo "$SHIMS_LINE" >> "$ZSHRC"
fi

asdf plugin add golang
asdf plugin add nodejs
asdf plugin add yarn

GOENV_LINE='. ${ASDF_DATA_DIR:-$HOME/.asdf}/plugins/golang/set-env.zsh'
if ! grep -Fxq "$GOENV_LINE" "$ZSHRC"; then
  echo "$GOENV_LINE" >> "$ZSHRC"
fi

chmod u+x '${ASDF_DATA_DIR:-$HOME/.asdf}/plugins/golang/set-env.zsh'

source "$ZSHRC"
