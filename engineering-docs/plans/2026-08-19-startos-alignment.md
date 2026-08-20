# Plano ativo: plataforma de serviços LabsOS

**Status:** posterior à fundação do appliance

## Objetivo

Evoluir o monólito LabsOS para uma plataforma doméstica de serviços self-hosted
para aplicações, APIs, bancos, automações, ferramentas de desenvolvimento,
agentes de IA e serviços pessoais. O StartOS é referência de produto e
lifecycle; Docker/Compose é o runtime inicial.

## Ordem de execução

1. Adicionar SQLite e extrair regras de domínio de `backend/internal/api` para um Core explícito.
2. Definir o contrato `ServicePackage` com metadata, dependências, interfaces,
   actions, tasks, health checks, migrações e backups.
3. Implementar `DockerRuntime` usando Docker Compose SDK.
4. Separar o contrato de serviço do adaptador Docker/Compose.
5. Transformar `runtime/` em protocolo real entre API, `labsd` e runtime.
6. Persistir eventos, jobs e estado com versionamento e replay.
7. Implementar backup/restore por serviço usando Restic.
8. Implementar updates e migrações por serviço.
9. Distribuir pacotes via OCI/ORAS e verificar com Cosign.
10. Avaliar isolamento adicional somente após estabilizar o modelo.

## Não fazer ainda

- Portar LXC: não é requisito do LabsOS.
- Criar um formato de pacote antes de definir o lifecycle.
- Mover mais diretórios apenas por semelhança visual com o StartOS.
- Declarar operações reais como prontas sem validação em QEMU/host de teste.

## Critério de saída

O plano só deve ser marcado como concluído quando houver um serviço de exemplo
instalável, configurável, monitorável, atualizável, sujeito a backup/restore e
removível pelo mesmo contrato, sem a UI depender de Compose diretamente.
