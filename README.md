# LabsOS

Base consolidada do LabsOS:

- `src/` e `public/`: frontend React atual;
- `backend/`: Labs API em Go;
- `specs/backend/`: contrato e arquitetura atuais do backend.

Home, Apps, App Store, Files e Settings usam a Labs API.

## Desenvolvimento

Frontend: `pnpm dev`.

## Desenvolvimento seguro

O backend Linux escuta apenas em `127.0.0.1` por padrão. Operações que podem
alterar disco, rede, reboot ou serviços ficam protegidas durante o desenvolvimento:

```bash
cd backend
LABSOS_ADDR=127.0.0.1:18080 go run ./cmd/labs-api
```

Para permitir acesso remoto e operações reais é necessário habilitar
explicitamente as proteções, apenas em uma máquina de teste:

```bash
LABSOS_ALLOW_REMOTE=true \
LABSOS_ENABLE_REAL_OPERATIONS=true \
LABSOS_CONFIRM_REAL_OPERATIONS=CONFIRMO \
LABSOS_ENABLE_NETWORK_CHANGES=true \
go run ./cmd/labs-api
```

Sem essas variáveis, o servidor permanece local e as operações perigosas não
devem produzir efeitos no host.

Em outro terminal, execute `pnpm dev`. O Vite encaminha `/api` para a API Linux local em `localhost:18080`, que coleta dados da máquina onde está rodando.

Detalhes: [`specs/backend/DEVELOPMENT.md`](specs/backend/DEVELOPMENT.md).
