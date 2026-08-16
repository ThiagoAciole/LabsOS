#!/usr/bin/env bash
set -Eeuo pipefail
root=$(cd "$(dirname "$0")/../.." && pwd); profile=${1:-production}; version=$(cat "$root/VERSION")
[[ "$(uname -s)" == Linux ]] || { echo 'ERRO: execute no WSL2/Linux.' >&2; exit 1; }
command -v simple-cdd >/dev/null || { echo 'ERRO: simple-cdd ausente. Rode make setup.' >&2; exit 1; }
mkdir -p "$root/build" "$root/dist"
for f in "$root/build/packages"/*.deb; do [[ -f "$f" ]] || { echo 'ERRO: rode make packages.' >&2; exit 1; }; done
cfg="$root/build/simple-cdd"; rm -rf "$cfg"; mkdir -p "$cfg/profiles/labsos"
cp "$root/iso/profiles/labsos/labsos.preseed" "$cfg/profiles/labsos/"
cp "$root/iso/profiles/labsos/package-lists" "$cfg/profiles/labsos/"
cp "$root/iso/profiles/labsos/labsos.conf" "$cfg/profiles/labsos/"
cp "$root/build/packages"/*.deb "$cfg/profiles/labsos/"
printf '\nlocal_packages="%s/build/packages"\nmirror_files=""\nexport OMIT_DOC_TOOLS=1\n' "$root" >> "$cfg/profiles/labsos/labsos.conf"
printf '%s\n' "$profile" > "$cfg/profiles/labsos/profile"
if command -v build-simple-cdd >/dev/null; then
  (cd "$cfg" && build-simple-cdd --force-root --dist trixie --profiles labsos --conf profiles/labsos/labsos.conf)
else
  echo 'ERRO: build-simple-cdd ausente.' >&2; exit 1
fi
iso=$(find "$cfg" -type f -iname '*.iso' | head -1)
[[ -n "$iso" ]] || { echo 'ERRO: simple-cdd não produziu ISO.' >&2; exit 1; }
cp "$iso" "$root/dist/LabsOS-$version-amd64.iso"; sha256sum "$root/dist/LabsOS-$version-amd64.iso" > "$root/dist/LabsOS-$version-amd64.iso.sha256"
echo "ISO: $root/dist/LabsOS-$version-amd64.iso"
