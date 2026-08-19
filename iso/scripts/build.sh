#!/usr/bin/env bash
set -Eeuo pipefail

root=$(cd "$(dirname "$0")/../.." && pwd)
profile=${1:-production}
version=$(cat "$root/VERSION")
[[ "$(uname -s)" == Linux ]] || { echo 'ERRO: execute em Ubuntu/Debian Linux ou WSL2.' >&2; exit 1; }
command -v lb >/dev/null || { echo 'ERRO: live-build ausente. Rode make setup.' >&2; exit 1; }

packages_dir="$root/build/packages"
[[ -d "$packages_dir" ]] || { echo 'ERRO: rode make packages.' >&2; exit 1; }
compgen -G "$packages_dir/*.deb" >/dev/null || { echo 'ERRO: nenhum pacote foi construído; rode make packages.' >&2; exit 1; }

build_dir="${LABSOS_IMAGE_BUILD_DIR:-$root/build/live-build}"
case "$build_dir" in
  /run/media/*|/media/*|/mnt/*)
    echo 'ERRO: LABSOS_IMAGE_BUILD_DIR precisa estar em filesystem Linux nativo (ex4/ext4), não em exFAT/NTFS.' >&2
    exit 1
    ;;
esac
if [[ -e "$build_dir" ]]; then
  if [[ "$EUID" -eq 0 ]]; then
    rm -rf -- "$build_dir"
  else
    sudo rm -rf -- "$build_dir"
  fi
fi
mkdir -p "$build_dir/config/includes.chroot/opt/labsos/packages" \
  "$build_dir/config/includes.chroot/etc/labsos" \
  "$build_dir/config/includes.chroot/usr/bin" \
  "$build_dir/config/includes.chroot/usr/lib/labsos" \
  "$build_dir/config/includes.chroot/etc/systemd/system" \
  "$build_dir/config/includes.chroot/root/isolinux" \
  "$build_dir/config/includes.binary" \
  "$build_dir/config/package-lists" \
  "$build_dir/config/bootloaders/isolinux" \
  "$build_dir/config/hooks/live" \
  "$build_dir/config/hooks/normal"

cp "$packages_dir"/*.deb "$build_dir/config/includes.chroot/opt/labsos/packages/"
cp "$root/iso/profiles/labsos/package-lists" "$build_dir/config/package-lists/labsos.list.chroot"
cp /usr/lib/ISOLINUX/isolinux.bin "$build_dir/config/includes.chroot/root/isolinux/isolinux.bin"
cp /usr/lib/syslinux/modules/bios/vesamenu.c32 "$build_dir/config/includes.chroot/root/isolinux/vesamenu.c32"
cp "$root/iso/packaging/labsos-core/rsvg-wrapper" "$build_dir/config/includes.chroot/usr/bin/rsvg"
chmod 0755 "$build_dir/config/includes.chroot/usr/bin/rsvg"
cp "$root/iso/packaging/labsos-core/labsos-installer" "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer"
chmod 0755 "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer"
cp "$root/iso/packaging/labsos-core/labsos-installer-backend" "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer-backend"
chmod 0755 "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer-backend"
cp /usr/share/live/build/bootloaders/isolinux/*.cfg "$build_dir/config/bootloaders/isolinux/"
cp /usr/share/live/build/bootloaders/isolinux/splash.svg.in "$build_dir/config/bootloaders/isolinux/"
cp /usr/lib/ISOLINUX/isolinux.bin "$build_dir/config/bootloaders/isolinux/isolinux.bin"
cp /usr/lib/syslinux/modules/bios/vesamenu.c32 "$build_dir/config/bootloaders/isolinux/vesamenu.c32"
printf '%s\n' "$version" > "$build_dir/config/includes.chroot/etc/labsos/image-version"
printf '%s\n' "$profile" > "$build_dir/config/includes.chroot/etc/labsos/image-profile"
touch "$build_dir/config/includes.binary/.labsos-live"

cat > "$build_dir/config/hooks/live/010-labs-user.hook.chroot" <<'EOF'
#!/bin/sh
set -eu
if ! id labs >/dev/null 2>&1; then
  useradd --create-home --shell /bin/bash --groups sudo labs
fi
# The live image needs a usable local account for LightDM and tty fallback.
# This is a temporary Live-session password; the installer must replace it.
printf '%s\n' 'labs:labs' | chpasswd
EOF
chmod 0755 "$build_dir/config/hooks/live/010-labs-user.hook.chroot"

cat > "$build_dir/config/hooks/normal/900-labsos-packages.hook.chroot" <<'EOF'
#!/bin/sh
set -eu
if ! id labs >/dev/null 2>&1; then
  useradd --create-home --shell /bin/bash --groups sudo labs
fi
printf '%s\n' 'labs:labs' | chpasswd
export DEBIAN_FRONTEND=noninteractive

dpkg -i /opt/labsos/packages/*.deb || apt-get -f install -y --allow-unauthenticated
rm -rf /opt/labsos/packages
EOF
chmod 0755 "$build_dir/config/hooks/normal/900-labsos-packages.hook.chroot"

cat > "$build_dir/config/includes.chroot/etc/hostname" <<'EOF'
labs
EOF

cat > "$build_dir/config/includes.chroot/etc/hosts" <<'EOF'
127.0.0.1 localhost labs
::1 localhost ip6-localhost ip6-loopback labs
EOF

cd "$build_dir"
lb config \
  --mode debian \
  --distribution trixie \
  --architectures amd64 \
  --binary-images iso-hybrid \
  --archive-areas "main contrib non-free non-free-firmware" \
  --linux-packages linux-image \
  --firmware-chroot false \
  --mirror-bootstrap https://deb.debian.org/debian/ \
  --mirror-chroot https://deb.debian.org/debian/ \
  --mirror-chroot-security https://security.debian.org/debian-security \
  --security false \
  --apt-secure false \
  --apt-options "--yes --allow-unauthenticated" \
  --apt-recommends false \
  --bootappend-live "boot=live components quiet splash username=labs user-password=labs" \
  --iso-application "LabsOS $version" \
  --iso-volume "LABSOS_$version" \
  --iso-preparer "LabsOS" \
  --iso-publisher "LabsOS"

# live-build 3.x não expõe LB_INITSYSTEM como opção de lb config.
printf '%s\n' 'LB_INITSYSTEM="systemd"' >> "$build_dir/config/common"

# live-build's syslinux stage expects a gfxboot cpio archive even when the
# optional gfxboot theme is not installed (apt recommends are disabled).
mkdir -p "$build_dir/binary/isolinux"
(cd "$build_dir/binary/isolinux" && printf '.' | cpio --quiet -o -H newc > bootlogo)
cp /usr/lib/ISOLINUX/isolinux.bin "$build_dir/binary/isolinux/isolinux.bin"
for module in ldlinux.c32 libcom32.c32 libutil.c32 libmenu.c32 vesamenu.c32; do
  cp "/usr/lib/syslinux/modules/bios/$module" "$build_dir/binary/isolinux/$module"
done
cp "$build_dir/config/bootloaders/isolinux/"*.cfg "$build_dir/binary/isolinux/"
cat > "$build_dir/binary/isolinux/live.cfg" <<'EOF'
label live
  menu label LabsOS Live
  linux /live/vmlinuz
  initrd /live/initrd.img
  append boot=live components quiet splash
EOF

if [[ "$EUID" -eq 0 ]]; then
  lb build
else
  sudo lb build
fi

iso=$(find "$build_dir" -maxdepth 1 -type f -name '*.iso' -print -quit)
[[ -n "$iso" ]] || { echo 'ERRO: live-build não produziu ISO.' >&2; exit 1; }
mkdir -p "$root/dist"
output="$root/dist/LabsOS-$version-amd64.iso"
cp "$iso" "$output"
sha256sum "$output" > "$output.sha256"
echo "ISO: $output"
