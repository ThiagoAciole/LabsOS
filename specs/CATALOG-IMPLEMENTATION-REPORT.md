# Relatório de implementação do catálogo

## Fonte e protocolo

- Fonte padrão Linux: `https://raw.githubusercontent.com/IceWhaleTech/CasaOS-AppStore/gh-pages/index.json`.
- Formato consumido: catálogo v2 publicado, com `apps[]`, `base_url` e metadata já resolvida.
- O repositório de origem usa `Apps/*/docker-compose.yml` e `x-casaos`; esses detalhes não são expostos pela Labs API.
- `LABSOS_CATALOG_URL` permite trocar a fonte em desenvolvimento/testes.

## Modelo e cache

- `backend/catalog` normaliza o payload para o modelo próprio do LabsOS.
- IDs de Jellyfin e Syncthing são canonizados para `jellyfin` e `syncthing`.
- Ícones relativos são resolvidos contra `base_url`.
- Cache: `/var/lib/labsos/catalog/apps.json`, gravação atômica e fallback para último cache válido.

## Apps e segurança

- `GET /api/v1/catalog/apps` representa descoberta.
- `GET /api/v1/apps` usa listagem real do Docker via operação tipada `ListApps` do `labsd`.
- Immich e Home Assistant não são retornados como instalados no Linux.
- Jellyfin permanece suportado.
- Syncthing possui manifest seguro e operações tipadas preparadas.
- Outros apps aparecem como `installable: false` / “Em breve”.
- A metadata editorial duplicada do frontend foi removida; a App Store usa apenas `GET /api/v1/catalog/apps`.
- Não há execução de Compose arbitrário, shell genérico, Docker socket, `privileged` ou mounts fora de `/DATA`.

## Validação

- `go test ./...`: aprovado.
- `go build ./...`: aprovado.
- `pnpm typecheck`: aprovado.
- `pnpm build`: aprovado anteriormente após a integração do catálogo.

## Pendências reais

- Deploy no Debian e validação do ciclo real do Syncthing.
- Reboot e smoke completo no host real.
- Validar a listagem `docker ps -a` no `labsd` no Debian.
