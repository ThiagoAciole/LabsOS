# Plano de entrega: ISO LabsOS operacional hoje

**Objetivo:** produzir hoje uma ISO LabsOS bootável em QEMU, com instalador visual em modo kiosk, instalação em disco virtual, first boot, dashboard React em kiosk e serviços Go ativos.

**Escopo de hoje:** MVP funcional e demonstrável. A ISO deve instalar e iniciar o sistema em QEMU sem depender de uma instalação manual do Debian. Compatibilidade ampla com todos os recursos do StartOS, hardware físico, criptografia, rollback completo e loja final de aplicativos ficam fora desta entrega.

## Resultado de aceite

Em uma máquina de teste limpa:

1. `make iso-dev` gera `dist/LabsOS-<version>-amd64.iso`.
2. A ISO inicializa em QEMU.
3. O sistema entra automaticamente no kiosk.
4. O instalador mostra discos detectados e oferece instalação no disco virtual.
5. A instalação cria `/DATA`, `/var/lib/labsos` e o estado de first boot.
6. Após reboot, `labs-api`, `labsd` e o dashboard iniciam via systemd.
7. O navegador abre o dashboard local em kiosk.
8. O dashboard consulta `/api/v1/system/summary`, `/api/v1/system/health` e `/api/v1/apps`.
9. O smoke test confirma serviços, API, dashboard, Docker e `/DATA`.
10. A ISO e seu checksum ficam em `dist/`.

## Limites explícitos

- Não implementar hoje suporte a S9PK.
- Não portar hoje o runtime LXC.
- Não implementar hoje o pipeline completo multi-hardware do StartOS.
- Não colocar credenciais padrão na imagem.
- Não executar operações destrutivas em disco real: os testes usam QEMU/qcow2.
- Não misturar esta entrega com o rebranding completo.
- Não apagar nem sobrescrever alterações já existentes no worktree; os arquivos modificados devem ser inspecionados antes de qualquer edição.

## Estado atual e riscos conhecidos

- O LabsOS já possui React/Vite, API Go, `labsd`, catálogo Compose, units systemd e um pipeline de imagem em migração para `live-build`.
- O pipeline direto `live-build` produz uma imagem live/hybrid; o instalador persistente e o first boot definitivo ainda são etapas seguintes.
- O serviço `labs-dashboard` serve arquivos estáticos; ainda é necessário garantir a sessão gráfica e o navegador kiosk.
- O repositório possui alterações locais não commitadas. Cada arquivo tocado deve ser revisado para evitar sobrescrever trabalho existente.
- O ambiente de build exige Debian/Ubuntu com `live-build`, `debootstrap`, `xorriso`, QEMU, Go, Node/npm e pnpm.

## Arquitetura do MVP

```text
ISO/USB
  ├── Debian live-build + rootfs preparado
  ├── pacote labsos-core
  ├── pacote labsos-api
  ├── pacote labsd
  ├── pacote labsos-web
  └── kiosk session

Boot inicial
  └── labsos-kiosk → React Installer
                         ↓ JSON-RPC/HTTP local
                       labs-installer
                         ↓
                      disco virtual

Boot instalado
  └── labsos-kiosk → React Dashboard
                         ↓
                       labs-api → labsd → Docker/systemd
```

## Plano de execução por blocos

### Bloco 0 — preservar e preparar o worktree

- Registrar `git status` e a lista dos arquivos já modificados.
- Não reformatar o repositório inteiro.
- Confirmar que `node_modules`, artefatos temporários e arquivos gerados não entram na ISO nem no commit.
- Confirmar `VERSION` e o nome final do artefato.

**Saída:** lista de arquivos em escopo e ambiente de build conhecido.

### Bloco 1 — fechar o build determinístico da ISO

Arquivos prováveis:

- `iso/scripts/build.sh`
- `iso/scripts/build-packages.sh`
- `iso/scripts/setup-build-env.sh`
- `iso/profiles/labsos/package-lists`
- `iso/profiles/labsos/labsos.preseed` (legado; não participa do build principal)
- `Makefile`

Tarefas:

- Garantir que a build do frontend execute `pnpm install --frozen-lockfile` apenas quando necessário.
- Compilar os três binários Go com versão embutida.
- Empacotar units e scripts no `.deb` correto.
- Instalar um display manager, servidor gráfico/compositor e navegador kiosk na imagem.
- Manter a instalação automática apenas no disco virtual de teste; não selecionar `/dev/sda` ou `/dev/nvme*` por padrão.
- Produzir ISO e `.sha256`.

**Verificação:** `make packages` e `make iso-dev`.

### Bloco 2 — kiosk de instalação e operação

Adicionar ou consolidar:

- `iso/packaging/labsos-core/labsos-kiosk.service`
- script de sessão kiosk;
- configuração de autologin;
- rota/flag de modo installer;
- rota de modo dashboard;
- fallback visual caso a API não responda.

Comportamento:

- sem instalação concluída: abrir `/installer`;
- com instalação concluída: abrir `/`;
- reiniciar o navegador se ele encerrar;
- esperar a API antes de abrir o dashboard;
- não iniciar kiosk duplicado;
- permitir modo de diagnóstico via variável/flag de boot.

**Verificação:** boot da ISO no QEMU e captura da tela; confirmar que o navegador abre sem interação manual.

### Bloco 3 — instalador Go mínimo e seguro

Criar um serviço de instalação separado, com acesso privilegiado controlado:

- `backend/cmd/labs-installer/`
- `backend/installer/`
- contrato em `backend/internal/api/` ou módulo RPC próprio.

Operações mínimas:

```text
installer.status
installer.disks
installer.inspect
installer.install
installer.progress
installer.reboot
```

Implementação do MVP:

- listar discos e tamanhos;
- aceitar somente o disco explicitamente selecionado;
- recusar discos do sistema atual e dispositivos removíveis por padrão;
- particionar o disco virtual de teste;
- criar filesystem;
- copiar/instalar o root filesystem preparado;
- criar `/DATA`;
- gravar `/etc/labsos/install-state`;
- habilitar first boot;
- reiniciar.

**Verificação:** instalar em qcow2 novo e confirmar que o sistema volta pelo disco instalado.

### Bloco 4 — first boot

Consolidar `labsos-firstboot` para:

- validar que a instalação terminou;
- criar diretórios persistentes;
- gerar segredo local;
- configurar hostname inicial;
- ativar Docker, `labs-api`, `labsd`, dashboard e kiosk;
- criar estado inicial do sistema;
- marcar `firstboot.done` somente após sucesso;
- deixar logs recuperáveis via journalctl.

**Verificação:** executar duas vezes e confirmar idempotência.

### Bloco 5 — transporte operacional e contratos

Para hoje, manter o HTTP local existente para reduzir risco, mas organizar os métodos como comandos tipados equivalentes ao modelo JSON-RPC do StartOS.

Contrato mínimo:

```text
system.summary
system.health
system.reboot
apps.list
apps.status
jobs.get
events.stream
```

Se o transporte JSON-RPC já estiver parcialmente implementado, usar JSON-RPC; caso contrário, não bloquear a ISO por uma migração completa do REST. O contrato deve ser separado do transporte para permitir a troca depois.

**Verificação:** testes de contrato Go/TypeScript e smoke test HTTP local.

### Bloco 6 — dashboard React em modo instalado

- Garantir que o build Vite produza `index.html` e assets em um diretório instalável.
- Adicionar uma tela explícita de estado do sistema: booting, healthy, degraded e recovery.
- Ligar o dashboard à API real, sem depender de mock.
- Manter Apps como catálogo Compose existente, sem tentar S9PK hoje.
- Garantir fallback quando Docker ou catálogo estiver indisponível.

**Verificação:** `pnpm typecheck`, `pnpm build` e acesso pelo navegador kiosk.

### Bloco 7 — validação automatizada

Executar na ordem:

```bash
pnpm typecheck
pnpm build
go test ./...
make packages
make iso-dev
make test
make smoke
```

O teste QEMU deve confirmar:

- ISO inicializa;
- instalação não escreve em disco real;
- reboot funciona;
- serviços sobem;
- API responde;
- dashboard responde;
- kiosk abre;
- `/DATA` existe;
- Docker está disponível.

## Critérios de bloqueio

Se algum item abaixo impedir a ISO hoje, reduzir o escopo na ordem indicada:

1. kiosk completo → usar dashboard servido localmente e teste por `curl`;
2. instalador visual completo → entregar primeiro a imagem live e o kiosk operacional;
3. runtime de apps → deixar catálogo em modo leitura;
4. JSON-RPC → manter API HTTP, mas separar contratos e jobs;
5. hardware real → validar somente QEMU.

Não reduzir:

- checksum da ISO;
- isolamento do teste em qcow2;
- ausência de senha padrão;
- idempotência do first boot;
- smoke test dos serviços.

## Entregáveis finais

- ISO em `dist/`;
- checksum SHA-256;
- script de teste QEMU;
- smoke test;
- documentação atualizada em `iso/README.md`;
- lista de limitações do MVP;
- relatório de comandos executados e resultados;
- plano separado para a segunda etapa: JSON-RPC completo, installer avançado, updates/rollback, backup e runtime de apps.

## Definição de pronto de hoje

O trabalho está pronto quando uma pessoa consegue:

1. executar `make iso-dev`;
2. iniciar a ISO no QEMU;
3. ver o kiosk;
4. instalar em um qcow2 novo;
5. reiniciar;
6. ver o dashboard React funcionando;
7. confirmar API, Docker e serviços pelo smoke test.

Isso entrega a fundação operacional inspirada no StartOS sem tentar concluir, no mesmo dia, a compatibilidade completa de runtime e aplicativos.
