# Provisioning do Docker no Debian

## Auditoria de 14/08/2026

- Debian 13.6, `x86_64`.
- Docker Engine, CLI, containerd, Buildx e Compose plugin estão instalados pelo repositório oficial Docker.
- Docker está habilitado e ativo; o daemon permanece no socket Unix local.
- `/DATA` foi criado no filesystem raiz existente.
- `agent` não tem sudo genérico; a política permite somente instalar/reiniciar a Labs API.
- A API permanece em `127.0.0.1:8080`.

## Decisão

O provisioning foi executado manualmente pelo administrador, sem reparticionar ou formatar o disco. `agent` continua sem sudo genérico e sem acesso ao grupo `docker`.

## Procedimento futuro

1. Definir explicitamente o dispositivo/filesystem de `/DATA`, sem formatar ou mover dados automaticamente.
2. Aprovar uma mudança de privilégio limitada para instalar o Docker pelo repositório oficial do Debian/Docker.
3. Instalar somente os pacotes necessários e habilitar a unit Docker.
4. Validar `docker version`, `docker info`, `systemctl is-active docker` e portas escutando.
5. Confirmar que o daemon não publica TCP em `2375/2376`.
6. Criar `/DATA` e a área de configuração do System File App somente após os passos anteriores. ✅

O runtime do LabsOS já trata Docker ausente como `available=false`; nenhum endpoint HTTP executa a CLI.
