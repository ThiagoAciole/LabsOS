# Estado do produto e gates

Este documento substitui a leitura dos planos históricos como fonte de status.

## Concluído ou operacional no código

- Monólito organizado na raiz com `backend`, `web`, `packaging` e `runtime`.
- API Go com providers Linux, Docker e mock.
- Catálogo remoto/cacheado e manifests Compose validados.
- Lifecycle básico de apps e jobs.
- Dashboard integrado à API em Home, Apps, Settings e recursos auxiliares.
- Auth, secrets, backups, notificações, auditoria, logs, rede e installer com
  contratos e testes de API.
- Pipeline de pacotes/ISO, kiosk e first-boot em evolução; o Live agora inclui
  diretamente autologin, sessão Kiosk e a rota `/installer` como entrada de
  instalação.
- Proteções de loopback, operações simuladas e validações de segurança.

## Parcial ou ainda não comprovado

- Instalação persistente completa, autologin/Kiosk e recovery em QEMU.
- Operações reais de reboot, rede, backup, restore, update e rollback.
- Runtime privilegiado separado e completo.
- Teste de dois apps reais ponta a ponta após boot da ISO.
- Validação física BIOS/UEFI e suporte a múltiplas arquiteturas.
- Lint frontend limpo; build e typecheck passam, mas há regras ESLint pendentes
  nos componentes UI.

## Lacunas de alinhamento com StartOS

- `ServicePackage` versionado em vez de Compose como unidade pública.
- Actions, tasks, health checks e dependências por serviço.
- Migrações e política de backup declaradas por serviço.
- Pacotes assinados e registry verificável.
- Estado persistente reativo e transporte de diffs.
- Isolamento por serviço além do Docker host.

## Direção de produto definida

O LabsOS não busca ser uma implementação completa do StartOS. A meta é um
appliance doméstico para aplicações, APIs, bancos, automações, ferramentas de
desenvolvimento, agentes de IA e serviços pessoais. Docker/Compose permanece
como runtime inicial; LXC é opcional.

Tecnologias escolhidas para acelerar a evolução:

- SQLite para estado persistente;
- Docker Compose SDK para o adaptador de runtime;
- Restic para backup por serviço;
- OCI/ORAS para distribuição de artefatos;
- Cosign para assinatura e verificação;
- `ServicePackage` como contrato independente dessas implementações.

## Prioridade imediata

Antes de ampliar o catálogo ou implementar o marketplace, a prioridade é fechar
o ciclo de appliance: instalar uma vez em disco persistente, preservar `/DATA`
e `/var/lib/labsos`, atualizar releases sem reinstalar a ISO e fazer rollback
quando o health check falhar. O plano correspondente está em
`engineering-docs/plans/2026-08-19-installable-updatable-foundation.md`.

## Fonte de verdade

- Estado atual: código e este arquivo.
- Contratos: `docs/backend/`.
- Segurança: `SECURITY-CHECKLIST.md`.
- Trabalho futuro: `engineering-docs/plans/`.
- Histórico: `engineering-docs/plans/archive/` e relatórios de implementação.
