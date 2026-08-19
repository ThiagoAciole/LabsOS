# Backend LabsOS

A Labs API é a única fronteira pública entre a interface e o sistema operacional. Ela expõe conceitos de produto por HTTP, delega comportamento a providers e não oferece execução arbitrária de comandos.

## Stack e estrutura

- Go 1.26 e `net/http`;
- `cmd/labs-api`: processo HTTP;
- `internal/api`: rotas, DTOs JSON e erros;
- `internal/platform`: modelos e contrato de provider;
- `providers/mock`: fixture em memória usada somente por testes;
- `providers/linux`: provider operacional que lê a máquina local e conversa com `labsd`.

O backend usa o provider Linux e coleta dados da máquina local. Operações sensíveis permanecem protegidas pela política de segurança e pelo `labsd`.

## Executar

```powershell
cd C:\Projetos\LabsOS\backend
go run ./cmd/labs-api
```

Servidor: `http://localhost:8080`.

## Índice

- [Arquitetura](ARCHITECTURE.md)
- [API](API.md)
- [Providers](PROVIDERS.md)
- [Desenvolvimento](DEVELOPMENT.md)
- [Roadmap](ROADMAP.md)
- [Relatório da migração](MIGRATION-REPORT.md)
