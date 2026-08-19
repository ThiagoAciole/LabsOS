#!/usr/bin/env bash
set -Eeuo pipefail
root=$(cd "$(dirname "$0")/../.." && pwd); iso=$(find "$root/dist" -maxdepth 1 -name 'LabsOS-*.iso' | sort | tail -1); [[ -f "$iso" ]] || { echo 'ISO não encontrada; rode make iso.' >&2; exit 1; }
command -v qemu-system-x86_64 >/dev/null || { echo 'qemu-system-x86_64 ausente; rode make setup.' >&2; exit 1; }
disk=${LABSOS_QEMU_DISK:-"$root/build/labsos-test.qcow2"}; mkdir -p "$(dirname "$disk")"; if [[ "${LABSOS_QEMU_RESET:-0}" == 1 ]]; then rm -f -- "$disk"; fi; [[ -f "$disk" ]] || qemu-img create -f qcow2 "$disk" 32G >/dev/null
exec qemu-system-x86_64 -enable-kvm -m "${LABSOS_QEMU_MEMORY:-2048}" -cpu max -drive "file=$disk,format=qcow2" -cdrom "$iso" -boot d -nic user,model=virtio
