# API HTTP atual

Base: `http://localhost:8080/api/v1`. Erros usam `{"error":{"code","message","details"}}`.

O frontend usa paths relativos `/api/v1/...`; em desenvolvimento, o proxy Vite aponta para `http://127.0.0.1:18080` sem alterar a porta padrão do produto.

| Método | Path | Finalidade | Request | Response | Erros | Status |
|---|---|---|---|---|---|---|
| GET | `/system/summary` | Resumo do sistema | — | `SystemSummary` | 503 no Linux incompleto | READY |
| GET | `/system/health` | Saúde agregada | — | `SystemHealth` | 503 no Linux incompleto | READY |
| POST | `/system/reboot` | Solicitar reboot | — | `Job` 202; simulado no mock | 400/503 | READY |
| POST | `/system/shutdown` | Solicitar shutdown | — | `Job` 202; simulado no mock | 400/503 | READY |
| GET | `/apps` | Produtos instalados | — | `App[]` | 503 | READY |
| GET | `/catalog/apps` | Catálogo disponível | — | `App[]` | 503 | READY |
| POST | `/apps/{id}/install` | Instalar produto | — | `Job` 202 | 400/404/503 | READY |
| POST | `/apps/{id}/start` | Iniciar produto | — | `Job` 202 | 400/404/503 | READY |
| POST | `/apps/{id}/stop` | Parar produto | — | `Job` 202 | 400/404/503 | READY |
| POST | `/apps/{id}/restart` | Reiniciar produto | — | `Job` 202 | 400/404/503 | READY |
| DELETE | `/apps/{id}` | Remover produto | — | `Job` 202 | 404/503 | READY |
| GET | `/settings` | Ler configurações | — | objeto JSON | 503 | READY |
| PUT | `/settings/system` | Alterar configurações do sistema | objeto JSON, até 1 MiB | objeto consolidado | 400/503 | READY |
| PUT | `/settings/network` | Alterar configurações de rede | objeto JSON, até 1 MiB | objeto consolidado | 400/503 | READY |
| GET | `/jobs/{id}` | Consultar job | — | `Job` | 404/503 | READY |
| GET | `/events` | Listar atividade | — | `Event[]` | 503 | READY |

## Contratos principais

`SystemSummary`: além de hostname, status, uptime, versão, CPU, memória, temperatura e storage, expõe `ipAddress`, `networkOnline`, `networkDownloadBytesPerSecond` e `networkUploadBytesPerSecond`. No WSL, temperatura pode ser `0` quando nenhuma thermal zone é exposta; o frontend apresenta o campo como indisponível.

`App`: `id`, `name`, `icon`, `description`, `status`, `version`, `url`, `updateAvailable`, `installed`. Estados produzidos atualmente: `running` e `stopped`; `installing` e `error` são PLANNED.

`Job`: `id`, `status`, `message`. Jobs são concluídos imediatamente e mantidos em memória no mock; execução assíncrona persistente é PARTIAL.

## Não migrado

Files e Storage são LEGACY nesta base: havia protótipos antigos, mas não foram promovidos a contrato atual. Autenticação, SSE e persistência são PLANNED.
