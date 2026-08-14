# LabsOS

Base consolidada do LabsOS:

- `src/` e `public/`: frontend React atual;
- `backend/`: Labs API em Go;
- `specs/backend/`: contrato e arquitetura atuais do backend.

Home, Apps, App Store e Settings usam a Labs API. Files permanece independente nesta etapa.

## Desenvolvimento

Frontend: `pnpm dev`.

Backend mock:

```powershell
cd C:\Projetos\LabsOS\backend
$env:LABSOS_MODE = "mock"
$env:LABSOS_ADDR = "127.0.0.1:18080"
go run ./cmd/labs-api
```

Em outro terminal, execute `pnpm dev`. O Vite encaminha `/api` para a API mock no Windows ou para o LinuxProvider no WSL, conforme o processo ativo em `localhost:18080`.

Detalhes: [`specs/backend/DEVELOPMENT.md`](specs/backend/DEVELOPMENT.md).
