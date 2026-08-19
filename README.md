# LabsOS

O LabsOS é um monólito de produto para transformar um servidor doméstico em um
appliance administrável pelo navegador. A raiz coordena todos os componentes
do produto, seguindo a ideia do `start-os` do StartOS.

## Estrutura do produto

```text
LabsOS/
├── backend/          # Core, API, providers e contratos operacionais
├── web/              # Dashboard, setup wizard e UI do appliance
├── runtime/          # Limite do runtime e agente privilegiado
├── packaging/        # ISO, pacotes Debian, kiosk e systemd
├── assets/           # Assets públicos do produto
├── docs/             # Contratos e documentação do produto
├── engineering-docs/ # Planos e registros de implementação
└── shared-libs/      # Contratos extraídos quando houver múltiplos consumidores
```

O frontend, backend, runtime e packaging pertencem ao mesmo produto e são
construídos juntos. A separação acima organiza responsabilidades sem criar
aplicações independentes.

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

Detalhes: [`docs/backend/DEVELOPMENT.md`](docs/backend/DEVELOPMENT.md).
