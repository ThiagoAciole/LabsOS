# Visão de produto do LabsOS

## Definição

O LabsOS é um sistema/appliance para um servidor doméstico local. Ele fornece
uma experiência simples para executar e administrar serviços pessoais
self-hosted pelo navegador.

## Serviços no escopo

O produto deve suportar, como primeira classe:

- aplicações self-hosted;
- APIs e webhooks;
- bancos de dados;
- automações e workers;
- ferramentas de desenvolvimento;
- agentes e serviços de IA;
- stacks com múltiplos containers;
- serviços pessoais de arquivos, mídia e comunicação.

O serviço pode ser um container único ou uma stack de vários containers. O
usuário não deve precisar entender essa diferença para instalar e operar o
serviço.

## Posição tecnológica

O Docker/Compose é o runtime inicial por compatibilidade com o ecossistema e
pela velocidade de entrega. Ele não é o contrato público do produto.

```text
ServicePackage
    ↓
DockerRuntime / Compose SDK
    ↓
Docker Engine
```

O StartOS é referência para experiência, lifecycle, saúde, tarefas, backups e
updates. O LabsOS não precisa copiar LXC, Rust, Patch-DB ou o formato `.s9pk`
para atingir seu objetivo.

## Promessa operacional

Para cada serviço, o LabsOS deve oferecer:

```text
descobrir → instalar → configurar → iniciar → usar
              ↓
        health, logs, jobs
              ↓
       backup → update → restore
```

## Não objetivos da primeira fase

- clusters e orquestração empresarial;
- compatibilidade obrigatória com LXC;
- substituir todas as ferramentas Linux;
- suporte a dezenas de usuários e RBAC empresarial;
- criar um formato proprietário de containers antes de validar o modelo de serviço.

## Decisões de arquitetura

- SQLite será a fonte persistente inicial do estado do LabsOS.
- Restic será o primeiro backend de backup por serviço.
- OCI/ORAS será a base de distribuição de artefatos.
- Cosign será usado para assinatura/verificação quando o registry for introduzido.
- O contrato `ServicePackage` será independente de Docker e do transporte HTTP.
