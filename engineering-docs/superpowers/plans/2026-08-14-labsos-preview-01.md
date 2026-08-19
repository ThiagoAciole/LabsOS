# Plano de implementação: LabsOS Preview 0.1

> **Para agentes executores:** executar nesta sessão, preservando alterações existentes e usando `sudo -n` via `agent` no servidor.

**Objetivo:** fechar o fluxo mínimo real `labs.local -> Dashboard -> App Store -> Compose -> lifecycle -> Files -> reboot`.

**Arquitetura:** manter Go/net/http, `labs-api` loopback e `labsd` como fronteira privilegiada. Generalizar o catálogo e o installer em torno de Compose, persistindo instalações em disco e consultando Docker para status.

**Stack:** Go, React/Vite/TypeScript, Docker Compose, systemd, Debian 13.

## Restrições globais

- Preservar `/DATA`, `root:data-admin`, setgid e dados existentes.
- Usar `SSH -> agent -> sudo -n`; não usar root SSH.
- Não adicionar installer por aplicativo.
- Não expor `labs-api` fora de `127.0.0.1:8080`.
- Validar com `go test ./...`, `go build ./...`, typecheck, build, smoke API e reboot real.

### Tarefa 1: Backend Compose e persistência

**Arquivos:** `backend/labsd/*`, `backend/providers/linux/provider.go`, `backend/providers/docker/*`, `backend/catalog/*`, `backend/apps/*`, `backend/internal/platform/provider.go`, testes correspondentes.

- Reutilizar contratos existentes.
- Adicionar fonte normalizada de Compose, validação de IDs/paths e persistência JSON em `/var/lib/labsos/apps`.
- Implementar `List`, `Install`, `Start`, `Stop`, `Restart`, `Remove`, `Status` e logs via `docker compose`/`docker inspect`.
- Manter Jellyfin/Syncthing compatíveis e permitir pelo menos dois catálogos reais pelo mesmo caminho.
- Testar parser/validação, persistência e lifecycle com runners injetáveis.

### Tarefa 2: API e frontend real

**Arquivos:** `backend/internal/api/server.go`, providers, `src/api/*`, `src/features/Apps/*`, `src/features/Settings/*`.

- Expor sources, detalhe, install/lifecycle/status/logs e update sem duplicar endpoints existentes.
- Remover mocks de Settings/update onde houver API real.
- Preservar estados loading, empty, offline, error e success.

### Tarefa 3: Release e deploy

**Arquivos:** `backend/deploy/*`, workflow de release existente ou novo, scripts e documentação mínima.

- Gerar dashboard e binários Linux versionados.
- Copiar para `/opt/labsos/releases/<version>`, atualizar `current` de forma reversível e reiniciar units.
- Garantir `labs.local`, Avahi/proxy, dashboard e units habilitados.

### Tarefa 4: Validação real

- Rodar testes locais e build.
- Fazer deploy, smoke test de Dashboard, Apps, Files e dois Apps distintos.
- Validar status/logs/remove e reboot real; só então declarar READY ou listar gate bloqueado.
