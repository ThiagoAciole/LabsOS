# Relatório de Implementação

## Estado atual

O checkout Windows continua sendo o frontend React/Vite e a Labs API Go. O Debian `labsos` continua com `labs-api.service` ativo e ouvindo somente em `127.0.0.1:8080`.

## Concluído nesta etapa

- Home: storage real via filesystem `/DATA`, com fallback para `/`.
- Runtime: `ContainerRuntime` detecta Docker sem expor daemon ou endpoint de execução.
- Apps: distinção `system`/`user` e lifecycle tipado para System Apps.
- Catálogo: modelo próprio, normalização, provider built-in, downloader remoto e cache em arquivo; no Debian atual a resposta validada ainda usa `source:"builtin"` porque nenhuma URL remota foi configurada.
- App Store: usa `category` e `source` da Labs API, mantendo metadata local somente como fallback.
- Segurança: validator bloqueia privileged, host network, Docker socket, traversal e mounts fora de `/DATA/Apps`.
- Apps Docker: contrato de produto e estados independentes de terminologia Docker.
- DockerAppsProvider: implementação fail-closed criada; sem Docker retorna indisponível e não instala nada.
- Jellyfin: manifest seguro criado, limitado a `/DATA/Apps/jellyfin` e `/DATA/Media` read-only; install, stop, start, restart, remove e install foram validados pela Labs API.
- Files: FileBrowser Quantum v1.5.1 provisionado e saudável em `127.0.0.1:8081`; frontend usa proxy local `/file-manager/`.

## Ainda pendente

- `labsd` está instalado como unit root no Debian, com socket `root:agent` em `/run/labsos/labsd.sock`.
- A vertical Jellyfin real está ativa e sua listagem/status já vem do `labsd`; aplicativos adicionais dependem de manifests declarativos válidos e do `labsd` local.
- Fonte remota ainda precisa de um endpoint compatível configurado e validado; não foi inventada uma URL externa.
- Reboot/recovery validado: SSH retornou, Docker/labsd/labs-api ficaram `active` e `enabled`, `/DATA` foi preservado, o socket foi recriado, Files respondeu `200` e Jellyfin respondeu `302`.
- Provisioning: Debian 13 `x86_64`, Docker Engine 29.7.2, Compose v5.4.0 e `/DATA` confirmados; sudo do `agent` permanece restrito.
- Segurança: checklist separado em `specs/SECURITY-CHECKLIST.md`, distinguindo evidências atuais dos gates ainda não comprovados.

## Validação

`go test ./...`, `go build ./...`, `pnpm typecheck` e `pnpm build` passaram. O build frontend mantém o aviso preexistente de bundle acima de 500 kB.

## Próximo passo

Configurar e validar uma fonte remota de catálogo compatível e adicionar manifests declarativos reais conforme cada aplicativo for validado individualmente.
