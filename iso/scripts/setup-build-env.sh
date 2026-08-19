#!/usr/bin/env bash
set -Eeuo pipefail
if [[ -z ${WSL_INTEROP:-} && ! -f /proc/version || "$(uname -s)" != Linux ]]; then
  echo 'ERRO: execute este script dentro do WSL2 (Debian/Ubuntu).' >&2; exit 1
fi
root=$(cd "$(dirname "$0")/../.." && pwd)
[[ "$root" == /mnt/* ]] && echo 'AVISO: em WSL2, o filesystem Linux costuma ser mais rápido que /mnt/*.'
missing=(); for p in live-build debootstrap xorriso isolinux syslinux-common syslinux-utils qemu-system-x86 qemu-utils make git curl wget jq rsync zstd dpkg fakeroot debhelper golang-go nodejs npm debian-archive-keyring gnupg; do dpkg -s "$p" >/dev/null 2>&1 || missing+=("$p"); done
if ((${#missing[@]})); then
  sudo apt-get update
  sudo apt-get install -y "${missing[@]}"
else echo 'Dependências do builder: OK'; fi
sudo npm install -g pnpm@10.15.0 >/dev/null
