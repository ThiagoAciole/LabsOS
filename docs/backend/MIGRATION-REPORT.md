# Relatório de migração

## Origem analisada

- `D:\labsos-draft\services\api`: módulo Go, servidor, interfaces e testes;
- `D:\labsos-draft\mocks`: árvore de filesystem mock;
- `D:\labsos-draft\openapi`: contrato conceitual antigo;
- `D:\labsos-draft\systemd`: units conceituais;
- `D:\labsos-draft\specs`: constituição, arquitetura, API, Apps, Settings, Security, Testing e Development Modes.

## Recuperado — KEEP

Tecnologia Go, conceitos de System/Apps/Settings, fixtures de teste, jobs/events e respostas JSON.

## Adaptado — ADAPT

O estado de teste foi isolado em provider de suporte; Apps agora expõem conceitos de produto; settings retornam cópias defensivas; ações geram jobs consultáveis; erros seguem envelope estruturado.

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

READY: system summary/health/power, apps/catalog/actions/remove, settings,
jobs/events, Files seguro em `/DATA`, Storage somente leitura, autenticação,
SSE, notificações, backups, rede, SSH, secrets, installer e exposição de serviços.
Operações privilegiadas continuam condicionadas à política local. Consulte `API.md`.

## Testes

Os testes cobrem transições de Apps, isolamento de Settings, energia protegida, rotas HTTP, jobs e erros estruturados. `go test ./...` e `go build ./...` passam no ambiente atual. O smoke HTTP valida health, summary, apps, settings, events, ações e consulta de jobs.

A porta 8080 pertence ao intervalo reservado `8072–8171` nesta máquina. O smoke usou `LABSOS_ADDR=127.0.0.1:18080`; o padrão do produto continua `127.0.0.1:8080`.

## Pendências operacionais

- validação em hardware dedicado das operações destrutivas e de rede;
- backend remoto de catálogo compatível configurado pelo administrador;
- OpenAPI regenerado a partir do contrato aprovado.

## Próxima etapa

Comparar os tipos/dados locais do frontend com `API.md`, definir mapeamentos e criar um client HTTP central por domínio. Não conectar antes dessa revisão de compatibilidade.
