# Plano de implementação: pipeline ISO LabsOS

> **Para agentes executores:** SUB-SKILL OBRIGATÓRIO: usar superpowers:executing-plans para implementar este plano tarefa por tarefa.

**Objetivo:** gerar uma ISO Debian Trixie amd64 reproduzível pelo WSL2.

**Arquitetura:** Makefile chama scripts Bash; Go e Vite são compilados em pacotes `.deb`; simple-cdd monta Debian Installer com preseed e pacotes LabsOS.

**Stack:** Debian 13, simple-cdd, dpkg-deb, Go, pnpm, Vite, QEMU.

## Restrições globais

- Nunca formatar discos físicos automaticamente.
- Não incluir segredos ou senhas na ISO.
- Preservar `labs-api`, `labsd` e `labs-dashboard` existentes.
- Produzir checksum SHA-256.

### Tarefa 1: builder e pacotes

- [x] Criar `VERSION`, `Makefile`, scripts de dependências e empacotamento.
- [x] Reutilizar os binários Go, units systemd e build Vite existentes.

### Tarefa 2: ISO e validação

- [x] Criar preseed, configuração simple-cdd, QEMU, smoke test, release e CI.
- [x] Validar dentro de WSL2 com `make setup` e `make packages`; `make iso` permanece bloqueado por build simple-cdd excedendo 4 minutos em `/mnt/c`.
