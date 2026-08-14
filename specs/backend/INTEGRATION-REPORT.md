# Relatório de integração

## Frontend analisado

Home, Apps, App Store, Settings e Files foram auditados. A UI aprovada foi preservada; apenas estados de loading, indisponibilidade, pending e erro foram conectados.

## Backend analisado

Foram usados System Summary/Health/Power, Apps/Catalog/Actions/Remove, Settings System/Network, Jobs e Events. O envelope `error.code/message/details` é interpretado no client central.

## Contratos alterados

`SystemSummary` ganhou IP, estado e tráfego de rede. DTOs de métricas, uptime, rede, apps e eventos são adaptados no frontend. O Vite encaminha `/api` para `localhost:18080` para suportar Windows e WSL2.

## Features conectadas

- **Home:** summary, health, CPU, RAM, uptime, temperatura disponível, rede, apps em execução e events vêm da API.
- **Apps:** listagem e drawer compartilham a mesma query; install/start/stop/restart/remove usam jobs síncronos do MockProvider.
- **App Store:** catálogo vem da API; categoria, origem, tamanho e destaques são metadata editorial local; instalação manual permanece indisponível.
- **Settings:** hostname, idioma, remoteAccess, summary, health, reboot e shutdown usam a API. Tema e updates permanecem locais.
- **Events:** Recent Activity usa `GET /events`; não há SSE nem polling complexo.

## Mocks removidos e restantes

System, rede, apps instalados, catálogo operacional, settings suportados e manutenção deixaram de ter mocks React como fonte de verdade. Permanecem metadata editorial da Store, updates, arquitetura, diagnóstico/relatório e relógio/clima. No modo Linux/WSL, Apps, Settings, Events e Jobs ainda usam fallback MockProvider.

## Pendências Backend

SQLite, jobs assíncronos, auth, SSE, update subsystem, `labsd`, Docker Apps e validação do SystemProvider em VM Debian.

## Files

Legacy Files não foi reintroduzido. A direção planejada continua sendo System File App / File Browser. Storage também não foi reintroduzido.

## Testes

Resultados desta execução devem ser lidos junto ao handoff: testes unitários do client/adapters, Go tests/build, frontend typecheck/build, lint específico e smoke HTTP mock foram executados. O lint global possui dívida preexistente nos primitives shadcn.

O LinuxProvider foi validado ao vivo no Ubuntu WSL2 e no servidor Debian 13.6. No Debian, a unit systemd ficou `enabled/active`, executando como `agent`, limitada a `127.0.0.1:8080`; o acesso do frontend ocorre por túnel SSH. Hostname, uptime, CPU, RAM, temperatura, IP e rede retornaram dados reais, e Power respondeu 503.

## Próximo passo

Implementar o provider real de Apps sem expor Docker na API pública. Temperatura continua indisponível apenas no WSL desta máquina.
