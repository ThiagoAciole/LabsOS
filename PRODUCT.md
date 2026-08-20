# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

O usuário principal mantém um home server ou homelab doméstico e quer uma experiência próxima de um produto pronto, sem administrar Linux, Docker ou serviços pelo terminal no uso diário.

Ele acessa o LabsOS pelo navegador, na rede local ou remotamente, em desktop, tablet ou celular.

## Product Purpose

O LabsOS permite transformar um servidor doméstico em uma plataforma simples
para executar aplicações self-hosted, APIs, bancos de dados, automações,
ferramentas de desenvolvimento, agentes de IA e outros serviços pessoais.

O usuário deve conseguir instalar, configurar, conectar, monitorar, atualizar e
fazer backup desses serviços sem administrar Linux, Docker ou Compose no uso
diário.

O fluxo mental principal é `Servidor -> LabsOS -> Apps / Files / Status`.

## Positioning

O LabsOS é appliance-first: apresenta Apps, Files e estado do servidor como conceitos do produto, enquanto esconde Docker, systemd, volumes, permissões, portas e a base Linux do fluxo comum.

## Operating Context

- Uso doméstico e homelab, não administração empresarial, clusters ou dezenas de usuários.
- Acesso local ou remoto pelo navegador.
- Apps e stacks são a unidade principal de descoberta, instalação e gerenciamento; containers são implementação interna.
- O runtime inicial é Docker/Compose por pragmatismo e compatibilidade com o ecossistema self-hosted.
- LXC não é requisito do produto; isolamento adicional só será considerado se os requisitos de segurança e operação exigirem.
- SSH pode existir para administração avançada, mas não é requisito para usar o produto.

## Capabilities and Constraints

- A App Store integrada é central para descobrir e instalar apps self-hosted.
- O Dashboard mostra somente informações úteis: CPU, RAM, disco, rede, uptime e apps.
- Files fica limitado a `/DATA` e não permite navegar livremente pelo filesystem do sistema.
- Arquivos do usuário, mídia, bancos e volumes persistentes ficam concentrados em `/DATA`.
- Operações destrutivas ou avançadas de armazenamento permanecem explícitas e manuais inicialmente.
- O frontend nunca acessa Docker ou Linux diretamente; toda operação passa pela API do LabsOS.
- A arquitetura atual pretendida é Debian Minimal, Labs Core/API em Go, Docker Engine, Docker Compose SDK, App Runtime, File Service e Labs Dashboard.
- O modelo de serviço deve evoluir para `ServicePackage`, usando Compose como adaptador interno inicial.
- O estado operacional deve ser persistido em SQLite e publicado por eventos para a interface.
- Backups devem usar uma estratégia por serviço, com Restic como primeira implementação.
- Distribuição futura deve usar OCI/ORAS e verificação de artefatos com Cosign.
- DietPi não é requisito arquitetural; foi apenas um ambiente de experimentação.

## Brand Commitments

- Nome do produto: LabsOS.
- A experiência deve parecer um sistema operacional ou appliance pronto, não um painel de administração Linux.
- A complexidade da infraestrutura deve permanecer invisível no fluxo comum.

## Evidence on Hand

- O repositório contém uma interface React/Vite com Dashboard, Apps, App Store, Files e Settings.
- O frontend já possui cliente HTTP configurável por `VITE_API_URL`.
- Os dados atuais de Dashboard e catálogo ainda incluem mocks e estados vazios; não há depoimentos, benchmarks ou outras provas externas que possam ser apresentados como fatos.

## Product Principles

1. Serviços pessoais, arquivos e saúde do servidor antes da infraestrutura.
2. Nenhum terminal obrigatório no fluxo comum.
3. Backend como única autoridade sobre o sistema.
4. Dados do usuário confinados a `/DATA`.
5. Operações destrutivas sempre explícitas.
6. Docker é implementação inicial; o contrato do produto é o serviço.
7. O StartOS é referência de experiência e lifecycle, não uma lista de tecnologias obrigatórias.

## Accessibility & Inclusion

A interface deve ser responsiva e utilizável em desktop, tablet e celular, incluindo acesso remoto pelo navegador.
