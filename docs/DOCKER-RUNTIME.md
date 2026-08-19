# Docker Runtime

O backend usa `backend/runtime.ContainerRuntime` como fronteira interna para o runtime de containers. A implementação atual detecta o Docker com `docker version --format` e expõe apenas disponibilidade e versão; handlers HTTP não executam comandos diretamente.

No Debian atual, Docker Engine 29.7.2 e Compose v5.4.0 estão instalados pelo repositório oficial para Debian Trixie. O daemon está habilitado, ativo e restrito ao socket Unix local; não há listener em 2375/2376.

Requisitos mantidos:

- `labs-api` continua em `127.0.0.1:8080`;
- Docker daemon não deve ser exposto via TCP;
- nenhum endpoint genérico de shell ou execução foi criado;
- operações privilegiadas futuras devem passar por comandos tipados, preferencialmente via `labsd`.

O provisioning criou `/DATA` no filesystem raiz existente, com grupo `data-admin` e permissões setgid 2770. O runtime da API continua sem acesso ao Docker daemon; operações reais de User Apps aguardam uma fronteira privilegiada tipada.

`labsd` é a fronteira tipada aprovada: roda como root, expõe somente `/run/labsos/labsd.sock` para o `agent` e aceita operações fixas para manifests conhecidos. O primeiro manifest executável é `docs/jellyfin.compose.yml`, com porta apenas em loopback e mounts limitados a `/DATA/Apps/jellyfin` e `/DATA/Media:ro`.

`providers/docker` já pode ser selecionado por uma camada futura de composição, mas permanece deliberadamente sem operações reais de instalação até haver runtime disponível e manifest validado.

O primeiro manifest preparado é `manifests.Jellyfin()`: configuração em `/DATA/Apps/jellyfin` e mídia em `/DATA/Media` somente leitura.
