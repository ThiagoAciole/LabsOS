# Relatório de migração

## Origem analisada

- `D:\labsos-draft\services\api`: módulo Go, servidor, interfaces e testes;
- `D:\labsos-draft\mocks`: árvore de filesystem mock;
- `D:\labsos-draft\openapi`: contrato conceitual antigo;
- `D:\labsos-draft\systemd`: units conceituais;
- `D:\labsos-draft\specs`: constituição, arquitetura, API, Apps, Settings, Security, Testing e Development Modes.

## Recuperado — KEEP

Tecnologia Go, modo mock no Windows, conceitos de System/Apps/Settings, fixtures de Apps, jobs/events e respostas JSON.

## Adaptado — ADAPT

O estado mock virou `MockProvider`; Apps agora expõem conceitos de produto; settings retornam cópias defensivas; ações geram jobs consultáveis; erros seguem envelope estruturado.

## Reescrito — REWRITE

O `main.go` monolítico e as interfaces não utilizadas foram substituídos por `cmd/labs-api`, `internal/api`, `internal/platform` e providers efetivamente injetados. O roteamento usa `net/http` nativo e o modo Linux falha fechado.

## Removido — DROP

Não foram transportados: frontend antigo, executável `labsos.exe`, Storage/Shares mock, Files mutável, OpenAPI desatualizado, scripts JS, installer e systemd units. Files antigo não protegia todo o contrato contra symlink escape; Storage expunha infraestrutura fora da prioridade atual.

## Estrutura final

```text
backend/
├── cmd/labs-api/main.go
├── internal/api/
├── internal/platform/
├── providers/mock/
├── providers/linux/
└── go.mod
```

## Contrato atual

READY: system summary/health/power, apps/catalog/actions/remove, settings system/network, jobs e events. PARTIAL: jobs síncronos em memória e Linux selecionável sem implementação. LEGACY: Files e Storage antigos. Consulte `API.md`.

## Testes

Os testes cobrem seleção de provider, transições de Apps, isolamento de Settings, energia simulada, rotas HTTP, jobs e erros estruturados. `go test ./...` e `go build ./...` passaram no Windows. O smoke HTTP validou health, summary, apps, settings, events, ação de App e consulta do job em modo mock.

A porta 8080 pertence ao intervalo reservado `8072–8171` nesta máquina. O smoke usou `LABSOS_ADDR=127.0.0.1:18080`; o padrão do produto continua `127.0.0.1:8080`.

## Pendências

- persistência SQLite e jobs realmente assíncronos;
- autenticação, autorização, CSRF/rate limiting e SSE;
- `labsd` e providers Linux;
- Files seguro em `/DATA`, se confirmado como próximo domínio;
- OpenAPI regenerado a partir do contrato aprovado.

## Próxima etapa

Comparar os tipos/dados locais do frontend com `API.md`, definir mapeamentos e criar um client HTTP central por domínio. Não conectar antes dessa revisão de compatibilidade.
