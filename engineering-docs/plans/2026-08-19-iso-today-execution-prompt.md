# Plano de execução hoje: ISO instalável do LabsOS

**Status:** prioritário

## Objetivo

Gerar hoje uma ISO LabsOS bootável em QEMU, instalar em um qcow2 novo, iniciar
o sistema instalado e comprovar que uma segunda release pode ser ativada sem
reinstalar a ISO.

## Escopo obrigatório

- ISO Debian/live-build bootável;
- installer seguro ou fluxo mínimo de instalação;
- instalação somente no disco explicitamente selecionado;
- criação e preservação de `/DATA` e `/var/lib/labsos`;
- first boot idempotente;
- `labs-api`, `labsd`, dashboard e Docker ativos;
- releases em `/opt/labsos/releases/<version>`;
- ponteiro `/opt/labsos/current`;
- update local simples sem reinstalar a ISO;
- smoke test e checksum.

## Fora do escopo de hoje

LXC, `ServicePackage` completo, marketplace, OCI/ORAS completo, Cosign,
Restic completo, suporte multi-arquitetura, hardware físico, networking
avançado e rollback complexo de dados.

## Ordem de execução

1. Diagnóstico e build: `git status`, `pnpm typecheck`, `pnpm build`,
   `go test ./...` e `bash -n packaging/scripts/*.sh`.
2. Pacotes e ISO: `make packages` e `make iso-dev`.
3. Teste QEMU: `make test` e `make smoke` usando qcow2 novo.
4. Update: instalar uma release ao lado da atual, validar, trocar `current`
   atomicamente, reiniciar, executar health check e manter rollback.
5. Relatório: registrar ISO, checksum, serviços ativos, gates e limitações.

Não usar `git reset`, `git checkout`, `rm -rf` ou limpeza ampla. Se uma etapa
secundária bloquear, preservar ISO bootável, installer seguro, first boot,
serviços básicos e checksum.

## Critério de pronto

```text
make iso-dev
→ boot ISO
→ instalar em qcow2
→ reboot
→ dashboard funcional
→ /DATA persistente
→ release nova ativada
→ update sem reinstalar ISO
```

Artefatos esperados:

```text
dist/LabsOS-<version>-amd64.iso
dist/LabsOS-<version>-amd64.iso.sha256
```
