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

# The Live session is the only place where the installer API may access the
# selected block device. This override is excluded by the installer backend
# before copying the system to the persistent disk.
mkdir -p "$build_dir/config/includes.chroot/etc/systemd/system/labs-api.service.d"
cat > "$build_dir/config/includes.chroot/etc/systemd/system/labs-api.service.d/live-installer.conf" <<'EOF'
[Service]
User=root
Group=root
PrivateDevices=false
Environment=LABSOS_ENABLE_REAL_OPERATIONS=true
Environment=LABSOS_CONFIRM_REAL_OPERATIONS=LIVE_SESSION
Environment=LABSOS_INSTALLER_ALLOW_DESTRUCTIVE=true
EOF

cp "$packages_dir"/*.deb "$build_dir/config/includes.chroot/opt/labsos/packages/"
cp "$root/packaging/profiles/labsos/package-lists" "$build_dir/config/package-lists/labsos.list.chroot"
cp /usr/lib/ISOLINUX/isolinux.bin "$build_dir/config/includes.chroot/root/isolinux/isolinux.bin"
cp /usr/lib/syslinux/modules/bios/vesamenu.c32 "$build_dir/config/includes.chroot/root/isolinux/vesamenu.c32"
cp "$root/packaging/labsos-core/rsvg-wrapper" "$build_dir/config/includes.chroot/usr/bin/rsvg"
chmod 0755 "$build_dir/config/includes.chroot/usr/bin/rsvg"
cp "$root/packaging/labsos-core/labsos-installer" "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer"
chmod 0755 "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer"
cp "$root/packaging/labsos-core/labsos-installer-backend" "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer-backend"
chmod 0755 "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-installer-backend"
# Install the LabsOS packages on the first Live boot. Some live-build
# versions do not execute custom chroot hooks when restoring a cached chroot;
# a systemd bootstrap is deterministic and runs before the kiosk.
cat > "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-live-bootstrap" <<'EOF'
#!/bin/sh
set -eu
if compgen -G '/opt/labsos/packages/*.deb' >/dev/null; then
  export DEBIAN_FRONTEND=noninteractive
  dpkg -i /opt/labsos/packages/*.deb || apt-get -f install -y --allow-unauthenticated
  rm -rf /opt/labsos/packages
fi
# Some live-build versions do not execute the live chroot hook that creates
# the local account. Create it here before LightDM attempts autologin.
if ! id labs >/dev/null 2>&1; then
  useradd --create-home --shell /bin/bash --groups sudo labs
fi
mkdir -p /DATA /DATA/Apps /DATA/Docker/Compose /DATA/Docker/Configs /DATA/Docker/Volumes /DATA/Media /DATA/Projects /DATA/Shared /DATA/Backups /var/lib/labsos/state
release=$(find /opt/labsos/releases -mindepth 1 -maxdepth 1 -type d | sort -V | tail -1)
if [ -n "$release" ]; then ln -sfn "$release" /opt/labsos/current; fi
# The live account must remain usable after all Debian package postinst
# scripts have run; this keeps the optional local account usable.
if id labs >/dev/null 2>&1; then printf '%s\n' 'labs:labs' | chpasswd; fi
systemctl daemon-reload
systemctl enable --now labs-api labsd labs-dashboard 2>/dev/null || true
EOF
chmod 0755 "$build_dir/config/includes.chroot/usr/lib/labsos/labsos-live-bootstrap"
cat > "$build_dir/config/includes.chroot/etc/systemd/system/labsos-live-bootstrap.service" <<'EOF'
[Unit]
Description=LabsOS Live package bootstrap
Before=labs-api.service labsd.service labs-dashboard.service
After=local-fs.target network-online.target
Wants=network-online.target

[Service]
Type=oneshot
ExecStart=/usr/lib/labsos/labsos-live-bootstrap
RemainAfterExit=yes

[Install]
WantedBy=multi-user.target
EOF
mkdir -p "$build_dir/config/includes.chroot/etc/systemd/system/multi-user.target.wants"
ln -sf ../labsos-live-bootstrap.service "$build_dir/config/includes.chroot/etc/systemd/system/multi-user.target.wants/labsos-live-bootstrap.service"
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
if ! id user >/dev/null 2>&1; then
  useradd --create-home --shell /bin/bash --groups sudo user
fi
printf '%s\n' 'labs:labs' | chpasswd
printf '%s\n' 'user:user' | chpasswd
export DEBIAN_FRONTEND=noninteractive

dpkg -i /opt/labsos/packages/*.deb || apt-get -f install -y --allow-unauthenticated
rm -rf /opt/labsos/packages
EOF
chmod 0755 "$build_dir/config/hooks/normal/900-labsos-packages.hook.chroot"

# live-build versions used by Debian can skip normal hooks when restoring a
# cached chroot. Repeat the package installation in the live stage, which is
# guaranteed to run against the final filesystem before squashfs is created.
cat > "$build_dir/config/hooks/live/900-labsos-packages.hook.chroot" <<'EOF'
#!/bin/sh
set -eu
if ! id labs >/dev/null 2>&1; then
  useradd --create-home --shell /bin/bash --groups sudo labs
fi
printf '%s\n' 'labs:labs' | chpasswd
export DEBIAN_FRONTEND=noninteractive
if compgen -G '/opt/labsos/packages/*.deb' >/dev/null; then
  dpkg -i /opt/labsos/packages/*.deb || apt-get -f install -y --allow-unauthenticated
  rm -rf /opt/labsos/packages
fi
EOF
chmod 0755 "$build_dir/config/hooks/live/900-labsos-packages.hook.chroot"

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
  --bootappend-live "boot=live components quiet splash systemd.ssh_auto=no" \
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

# live-build emits versioned kernel/initrd names on this Debian release, while
# our custom Syslinux entry uses stable names. Point the entry at the actual
# files and rebuild only the binary ISO stage.
kernel=$(find "$build_dir/binary/live" -maxdepth 1 -type f -name 'vmlinuz-*' -printf '%f\n' | sort | head -1)
initrd=$(find "$build_dir/binary/live" -maxdepth 1 -type f -name 'initrd.img-*' -printf '%f\n' | sort | head -1)
[[ -n "$kernel" && -n "$initrd" ]] || { echo 'ERRO: kernel/initrd não encontrados na imagem live.' >&2; exit 1; }
# Keep stable paths for Syslinux and for older tooling while retaining the
# versioned files emitted by live-build.
if [[ "$EUID" -eq 0 ]]; then
  ln -sfn "$kernel" "$build_dir/binary/live/vmlinuz"
  ln -sfn "$initrd" "$build_dir/binary/live/initrd.img"
else
  sudo ln -sfn "$kernel" "$build_dir/binary/live/vmlinuz"
  sudo ln -sfn "$initrd" "$build_dir/binary/live/initrd.img"
fi
sed -i \
  -e "s|linux /live/vmlinuz$|linux /live/$kernel|" \
  -e "s|initrd /live/initrd.img$|initrd /live/$initrd|" \
  "$build_dir/binary/isolinux/live.cfg"
source_iso=$(find "$build_dir" -maxdepth 1 -type f -name '*.iso' -print -quit)
[[ -n "$source_iso" ]] || { echo 'ERRO: live-build não produziu ISO.' >&2; exit 1; }
iso="$build_dir/LabsOS-patched.iso"
xorriso \
  -indev "$source_iso" \
  -outdev "$iso" \
  -boot_image any replay \
  -rm /isolinux/live.cfg -- \
  -map "$build_dir/binary/isolinux/live.cfg" /isolinux/live.cfg \
  -commit >/dev/null 2>&1
[[ -f "$iso" ]] || { echo 'ERRO: não foi possível corrigir o live.cfg na ISO.' >&2; exit 1; }
mkdir -p "$root/dist"
output="$root/dist/LabsOS-$version-amd64.iso"
cp "$iso" "$output"
sha256sum "$output" > "$output.sha256"
echo "ISO: $output"
