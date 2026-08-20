# Arquitetura atual do LabsOS

O LabsOS é um monólito de produto para um servidor pessoal. Web, Core/API,
runtime e packaging são componentes do mesmo appliance e são construídos a
partir desta raiz.

```text
┌─────────────────────── LabsOS appliance ───────────────────────┐
│  web/ → HTTP/SSE → backend/ → providers + labsd → Docker/Linux │
│                         │                                      │
│                     packaging/                                │
└────────────────────────────────────────────────────────────────┘
```

## Responsabilidades

- `web/`: dashboard, setup/installer e estados operacionais.
- `backend/internal/api/`: transporte HTTP e adaptação de requests/responses.
- `backend/internal/platform/`: contratos entre API e providers.
- `backend/providers/`: implementações Linux, Docker e mocks de teste.
- `backend/labsd/`: fronteira atual para operações privilegiadas tipadas.
- `backend/catalog/`, `apps/` e `manifests/`: descoberta, validação e lifecycle
  inicial de apps Compose.
- `packaging/`: imagem live, pacotes Debian, kiosk e unidades systemd.
- `runtime/`: limite arquitetural reservado para o runtime de serviços; a
  implementação atual ainda está em `backend/labsd` e `backend/runtime`.
- `shared-libs/`: contratos extraídos somente quando houver múltiplos
  consumidores reais.

## Estado arquitetural

O transporte atual é HTTP/JSON com stream de eventos. O frontend usa queries,
mutations e invalidação por eventos. Não há ainda um store reativo central
equivalente ao Patch-DB do StartOS.

Apps são atualmente descritos por manifests/Compose. O objetivo é introduzir
`ServicePackage` como contrato público, deixando Docker Compose como adaptador
interno. O pacote deverá declarar ações, tarefas, health checks, dependências,
migrações e política de backup.

O runtime inicial será Docker/Compose SDK. LXC é opcional e não faz parte do
escopo obrigatório do produto.

## Princípios

1. O frontend nunca acessa Docker, Linux ou filesystem diretamente.
2. Providers são a fronteira de ambiente e devem falhar fechado.
3. Operações privilegiadas passam por comandos tipados e auditáveis.
4. Dados do usuário ficam limitados às árvores aprovadas, especialmente `/DATA`.
5. O contrato público deve falar em apps/serviços, não em detalhes de containers.
