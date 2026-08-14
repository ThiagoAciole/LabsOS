# Backend LabsOS

A Labs API é a única fronteira pública entre a interface e o sistema operacional. Ela expõe conceitos de produto por HTTP, delega comportamento a providers e não oferece execução arbitrária de comandos.

## Stack e estrutura

- Go 1.26 e `net/http`;
- `cmd/labs-api`: processo HTTP;
- `internal/api`: rotas, DTOs JSON e erros;
- `internal/platform`: modelos e contrato de provider;
- `providers/mock`: estado funcional em memória para Windows;
- `providers/linux`: limite explícito para a futura integração com `labsd`.

## Modos

`LABSOS_MODE=mock` é funcional em Windows, macOS e Linux. `LABSOS_MODE=linux` seleciona o provider real, mas suas operações retornam `503 PROVIDER_UNAVAILABLE` até a implementação do `labsd`.

## Executar

```powershell
cd C:\Projetos\LabsOS\backend
$env:LABSOS_MODE = "mock"
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
