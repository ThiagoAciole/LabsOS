# Plano rápido: LabsOS como plataforma de serviços

**Status:** posterior à fundação do appliance

Este plano começa somente depois de concluir o plano prioritário de instalação
persistente e updates sem reinstalar a ISO.

## Objetivo

Entregar rapidamente um vertical slice completo do LabsOS: um serviço
`ServicePackage` que possa ser descoberto, instalado, configurado, iniciado,
monitorado, submetido a backup, atualizado e removido pela UI, usando Docker
Compose SDK, SQLite, Restic, OCI/ORAS e Cosign.

O objetivo não é implementar toda a paridade com StartOS de uma vez. É provar
o lifecycle completo com uma arquitetura que possa ser repetida para os demais
serviços.

## Decisões fixas

- Docker/Compose é o runtime inicial.
- LXC fica fora do escopo.
- SQLite é a fonte de estado do Core.
- Restic é o provider inicial de backup.
- OCI/ORAS é o transporte de pacotes.
- Cosign é a verificação de artefatos.
- HTTP/JSON local continua sendo o transporte inicial de `labsd`.
- O frontend continua React/Vite e a API continua Go.

## Serviço piloto

Usar um serviço simples e representativo, preferencialmente Jellyfin ou
Syncthing, com:

- uma imagem principal;
- volume persistente;
- interface web;
- health check;
- configuração mínima;
- backup restaurável.

Não começar com uma stack que exija migração complexa de banco. PostgreSQL e
agentes de IA entram depois que o lifecycle básico estiver comprovado.

## Fase 0 — preparação e contratos (meio dia)

### Entregáveis

- atualizar `docs/STATUS.md` com o serviço piloto;
- criar `backend/core/` sem mover tudo de uma vez;
- definir interfaces Go para `ServicePackage`, `ServiceRuntime`,
  `StateStore`, `BackupProvider` e `ArtifactRegistry`;
- definir estados comuns:

```text
discovered → installing → stopped → starting → running
                                     └──────→ degraded
running → stopping → stopped → removing → removed
```

### Regra

Nenhum handler novo deve receber regra de lifecycle. Ele chama o Core.

## Fase 1 — SQLite e Core mínimo (1 dia)

### Implementar

```text
backend/core/
├── state/
├── services/
├── jobs/
└── events/
```

Persistir:

- serviços instalados;
- versão;
- estado;
- jobs;
- eventos;
- notificações.

### Critério de aceite

- reiniciar a API não perde serviços nem jobs;
- `apps.list` usa SQLite;
- ações criam jobs persistidos;
- eventos possuem ID e sequência;
- testes de migração do banco passam.

## Fase 2 — `ServicePackage` e adapter Compose (1–2 dias)

### Implementar

```text
backend/core/services/
├── package.go
├── manifest.go
├── lifecycle.go
├── health.go
└── actions.go

backend/providers/docker/
└── compose_runtime.go
```

O parser atual de Compose deve gerar um `ServicePackage` provisório.

### Contrato mínimo

- ID e versão;
- imagens;
- volumes;
- interfaces/portas;
- dependências;
- actions;
- health checks;
- runtime;
- backup policy.

### Critério de aceite

- Jellyfin/Syncthing pode ser carregado como `ServicePackage`;
- a UI não precisa ler Compose diretamente;
- o runtime pode ser substituído por outro adapter sem alterar o Core.

## Fase 3 — lifecycle real com Compose SDK (1–2 dias)

Implementar, nesta ordem:

```text
install → start → health → logs → stop → restart → remove
```

Cada operação deve:

1. criar um job no SQLite;
2. publicar progresso;
3. executar o adapter Docker;
4. consultar o resultado real;
5. atualizar o estado;
6. registrar auditoria.

O Docker Compose SDK deve substituir gradualmente parsing de stdout e chamadas
manuais à CLI. O contrato atual pode permanecer compatível durante a migração.

### Critério de aceite

- instalar e remover o serviço piloto;
- reiniciar a API durante um job sem corromper o estado;
- health check diferencia `running`, `degraded` e `failed`;
- logs e erros aparecem na UI.

## Fase 4 — estado reativo (meio dia)

Manter SSE inicialmente:

```text
Core persiste evento
→ EventHub publica evento
→ SSE entrega evento
→ frontend atualiza query/estado
```

Padronizar o evento:

```json
{
  "id": "evt_123",
  "sequence": 42,
  "type": "service.updated",
  "aggregate": "service",
  "aggregateId": "jellyfin",
  "payload": {}
}
```

WebSocket só deve ser adotado se SSE não atender o dashboard.

## Fase 5 — backup e restore com Restic (1–2 dias)

Adicionar backup apenas ao serviço piloto.

### Fluxo

```text
stop
→ pre-backup
→ Restic snapshot
→ verify
→ metadata SQLite
→ start
→ health
```

### Critério de aceite

- backup criptografado criado;
- snapshot aparece na UI;
- restore em diretório/volume de teste;
- verificação de integridade;
- falha de restore não apaga o estado original.

## Fase 6 — update e migration (1 dia)

Implementar update do serviço piloto sem migração complexa:

```text
backup
→ stop
→ atualizar artefato
→ start
→ health
→ confirmar
```

Adicionar o contrato de migration mesmo que a primeira implementação seja
`no-op`.

### Critério de aceite

- update preserva volume;
- versão fica registrada no SQLite;
- falha pós-update marca `degraded`;
- rollback para a imagem anterior funciona.

## Fase 7 — OCI/ORAS e Cosign (1–2 dias)

Primeiro empacotar apenas o serviço piloto.

### Artefato

```text
ServicePackage
├── manifest.json
├── health.json
├── backup.json
├── actions.json
└── compose.yaml
```

Fluxo:

```text
build
→ OCI artifact
→ push registry
→ Cosign sign
→ pull
→ Cosign verify
→ install
```

Durante o desenvolvimento, permitir um registry local e sideload assinado.
Não criar ainda marketplace público.

## Fase 8 — UI e expansão (1–2 dias)

Atualizar a UI para usar o lifecycle do Core:

- tela de detalhes do serviço;
- status e health;
- actions declaradas pelo pacote;
- tasks pendentes;
- jobs e progresso;
- backup/restore;
- update/rollback;
- logs.

Depois migrar o segundo serviço e verificar que nenhum código específico do
piloto contaminou o Core.

## Paralelização segura

Podem ocorrer em paralelo:

```text
Core/SQLite       ─┐
ServicePackage    ─┼─→ lifecycle piloto
UI de lifecycle   ─┘

Restic provider   ─→ backup/restore
OCI/Cosign        ─→ distribuição do pacote
```

Não paralelizar ainda:

- LXC;
- networking avançado;
- marketplace público;
- múltiplos bancos complexos;
- refatoração completa de todos os handlers;
- migração para WebSocket/Patch-DB.

## Fase posterior

Depois do vertical slice:

1. migrar todos os apps para `ServicePackage`;
2. adicionar PostgreSQL/MySQL dump strategies;
3. adicionar tasks interativas e dependências;
4. registry remoto e publishers;
5. proxy/mDNS/TLS por interface;
6. installer/recovery completo;
7. limites de CPU/memória e isolamento mais forte;
8. avaliar Incus/LXC apenas se necessário.

## Definition of Done

O plano rápido está concluído quando o serviço piloto consegue:

```text
catalogar
→ verificar assinatura
→ instalar
→ configurar
→ iniciar
→ reportar health/logs
→ executar action
→ fazer backup
→ restaurar
→ atualizar
→ fazer rollback
→ remover
```

Tudo deve ser operável pela UI, persistido no SQLite, executado pelo runtime
Docker e coberto por testes de contrato e um smoke test de QEMU/ambiente Linux.
