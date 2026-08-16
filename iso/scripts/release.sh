#!/usr/bin/env bash
set -Eeuo pipefail
root=$(cd "$(dirname "$0")/../.." && pwd); version=$(cat "$root/VERSION"); iso="$root/dist/LabsOS-$version-amd64.iso"; [[ -f "$iso" ]] || { echo 'ISO ausente.' >&2; exit 1; }; sha256sum -c "$iso.sha256"
