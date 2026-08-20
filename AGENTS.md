# AGENTS.md — LabsOS

## Objetivo

O LabsOS é um monólito de produto para um servidor doméstico que executa
aplicações self-hosted, APIs, bancos de dados, automações, ferramentas de
desenvolvimento, agentes de IA e serviços pessoais simples. A organização e o
lifecycle são inspirados no StartOS, mas Docker/Compose é o runtime inicial.

## Estrutura do repositório

```text
backend/          API Go, Core atual, providers e comandos
web/              Dashboard React, setup e installer web
runtime/          Limite arquitetural do runtime do appliance
packaging/        ISO, pacotes Debian, kiosk e systemd
assets/           Assets públicos do produto
docs/             Arquitetura, contratos, segurança e relatórios
engineering-docs/ Planos ativos e histórico de engenharia
shared-libs/      Contratos compartilhados entre múltiplos consumidores
```

Não crie aplicações independentes dentro do repositório sem uma decisão
explícita. A separação de pastas representa domínios do mesmo appliance.

## Regras arquiteturais

- O frontend só conversa com a Labs API; nunca acessa Docker, Linux, `/proc`,
  `/sys`, discos ou filesystem real diretamente.
- `backend/internal/api` deve cuidar de transporte HTTP, DTOs e erros. Regras
  de domínio novas devem ser extraídas para um Core/domínio antes de aumentar
  handlers já grandes.
- Providers são a fronteira para dependências do ambiente e devem falhar
  fechado quando uma capacidade não estiver disponível.
- Operações privilegiadas passam por comandos tipados, auditáveis e limitados.
- Não adicione endpoints `/shell`, `/exec`, `/command` ou concatenação de input
  do usuário em `sh -c`.
- O contrato público deve falar em apps/serviços, não em containers, volumes ou
  detalhes internos de Docker.
- `runtime/` é a fronteira arquitetural futura; a implementação atual ainda
  está em `backend/labsd` e `backend/runtime`. Não mova código apenas por
  simetria visual: preserve o módulo Go e os contratos até haver necessidade
  real.
- `shared-libs/` só deve receber código usado por mais de um consumidor.

## Segurança operacional

Por padrão, a API deve permanecer em `127.0.0.1`. Operações de disco, rede,
reboot, shutdown, update, rollback, secrets e serviços devem continuar
simuladas ou protegidas durante desenvolvimento.

Para operações reais em uma máquina de teste, todas as confirmações precisam
ser explícitas:

```bash
LABSOS_ALLOW_REMOTE=true \
LABSOS_ENABLE_REAL_OPERATIONS=true \
LABSOS_CONFIRM_REAL_OPERATIONS=CONFIRMO \
LABSOS_ENABLE_NETWORK_CHANGES=true \
go run ./cmd/labs-api
```

Nunca execute testes destrutivos em disco físico. Use QEMU/qcow2 e verifique o
dispositivo alvo antes de qualquer operação de instalação.

## Documentação

- Estado atual: `docs/STATUS.md`.
- Arquitetura vigente: `docs/ARCHITECTURE.md`.
- API e contratos: `docs/backend/`.
- Segurança: `docs/SECURITY-CHECKLIST.md`.
- Planos ativos: `engineering-docs/plans/`.
- Planos históricos: `engineering-docs/plans/archive/`.

Não use planos arquivados como instruções atuais. Ao concluir ou invalidar um
plano, atualize `docs/STATUS.md` e o índice dos planos.

## Comandos de desenvolvimento

Frontend, a partir da raiz:

```bash
pnpm install --frozen-lockfile
pnpm dev
pnpm typecheck
pnpm build
```

Backend:

```bash
cd backend
go test ./...
go build ./...
```

Validação de scripts:

```bash
bash -n packaging/scripts/*.sh
git diff --check
```

Pipeline de imagem, em ambiente Debian/Ubuntu apropriado:

```bash
make packages
make iso-dev
make test
make smoke
```

O lint frontend pode conter falhas conhecidas nos componentes UI gerados. Não
declare a qualidade concluída sem registrar o resultado de `pnpm lint`.

## Modelo de serviço e StartOS

O LabsOS usa Compose/manifests como etapa atual. A evolução é um
`ServicePackage` com metadata, dependências, interfaces, actions, tasks, health
checks, migrações e política de backup. O pacote é o contrato; Docker Compose é
somente o adaptador inicial.

Tecnologias adotadas para acelerar o produto:

- SQLite para estado persistente;
- Docker Compose SDK para o runtime inicial;
- Restic para backup/restore;
- OCI/ORAS para distribuição;
- Cosign para assinatura e verificação.

Ao implementar essa evolução:

1. defina o contrato antes do adaptador Docker;
2. mantenha o lifecycle independente do transporte HTTP;
3. associe backup, restore e migração ao serviço;
4. exija health check após instalação e atualização;
5. não trate LXC como requisito; só avalie isolamento adicional depois que o
   modelo de serviço estiver estável.

## Alterações de arquivos

- Preserve alterações existentes no working tree.
- Não reformatte o repositório inteiro por causa de uma mudança localizada.
- Use `apply_patch` para edições locais e `git mv` para reorganizações de
  arquivos; mantenha renomeações rastreáveis.
- Não inclua `build/`, `dist/`, `*.iso`, discos virtuais, `.tmp-runtime/`,
  `vendor/` ou anexos locais em commits.
- Antes de apagar algo, confirme que não é código, documentação normativa,
  referência do StartOS ou saída necessária para uma validação em andamento.
- Commits devem agrupar uma mudança coerente e descrever o resultado, por
  exemplo: `feat: add service package contract` ou `refactor: separate core domain`.

## Critério de entrega

Antes de concluir uma alteração relevante:

1. atualize a documentação correspondente;
2. rode os testes do domínio afetado;
3. rode `pnpm typecheck` e `pnpm build` quando tocar no frontend ou config;
4. rode `go test ./...` quando tocar no backend;
5. execute `git diff --check`;
6. informe claramente validações que não puderam ser executadas.
