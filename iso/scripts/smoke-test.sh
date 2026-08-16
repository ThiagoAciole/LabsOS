#!/usr/bin/env bash
set -Eeuo pipefail
fail=0; for unit in docker labs-api labsd avahi-daemon; do systemctl is-active --quiet "$unit" || { echo "FALHA: $unit"; fail=1; }; done
[[ -d /DATA ]] || { echo 'FALHA: /DATA'; fail=1; }; docker version >/dev/null 2>&1 || fail=1; curl -fsS http://127.0.0.1:8080/ >/dev/null 2>&1 || echo 'AVISO: endpoint raiz da API não respondeu'; curl -fsS http://127.0.0.1/ >/dev/null 2>&1 || echo 'AVISO: dashboard não respondeu'; exit "$fail"
