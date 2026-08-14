# Checklist de Segurança

## Evidência atual

- [x] Labs API padrão em `127.0.0.1:8080`.
- [x] Debian validado via túnel SSH; nenhum bind LAN da API.
- [x] Nenhum endpoint `/shell`, `/exec` ou `/command`.
- [x] Docker daemon não instalado nem exposto.
- [x] `ContainerRuntime` é a única fronteira de CLI Docker.
- [x] `DockerAppsProvider` falha fechado sem runtime.
- [x] Manifests bloqueiam `privileged`, host network, Docker socket, traversal e mounts fora das árvores permitidas.
- [x] File Manager usa somente `/DATA` no manifesto de referência.
- [x] Configuração do File Manager fica fora dos dados do usuário.
- [x] Nenhuma formatação, movimentação ou exclusão de dados foi executada.

## Gates ainda não comprovados

- [ ] Docker instalado com política de privilégio aprovada.
- [ ] `/DATA` criado sem escolher ou formatar disco automaticamente.
- [ ] File Manager saudável após boot.
- [ ] Jellyfin instalado somente por manifest validado.
- [ ] Reboot e recovery validados.
