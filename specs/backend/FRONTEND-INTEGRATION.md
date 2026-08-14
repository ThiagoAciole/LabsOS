# Integração Frontend ↔ Backend

| UI | Necessidade | Endpoint | Backend DTO | Frontend model / transformação | Estado |
|---|---|---|---|---|---|
| Home / saúde | estado agregado | `GET /system/health` | `SystemHealth` | texto de produto por `status` e `mode` | ADAPT |
| Home / CPU, RAM, disco | métricas atuais | `GET /system/summary` | `SystemSummary` | bytes para GiB e percentuais; uptime em dias/horas/minutos | ADAPT |
| Home / apps | apps em execução | `GET /apps` | `App[]` | filtro `installed && status === running` | MATCH |
| Home / atividade | eventos recentes | `GET /events` | `Event[]` | ícone e texto derivados de `type`; API não informa horário | ADAPT |
| Home / rede | tráfego, IP e conectividade | `GET /system/summary` | `SystemSummary` | bytes/s para MB/s e fallback explícito | ADAPT |
| Apps / instalados | catálogo instalado | `GET /apps` | `App[]` | adapter remove detalhes de infraestrutura | ADAPT |
| Apps / ações | instalar, iniciar, parar, reiniciar, remover | `POST /apps/{id}/{action}`, `DELETE /apps/{id}` | `Job` | aguarda resposta síncrona e invalida queries | ADAPT |
| App Store | catálogo | `GET /catalog/apps` | `App[]` | categoria, origem, tamanho e destaques seguem metadata local | ADAPT |
| App Store / manual | instalação por manifesto | — | — | botão preservado e indisponível | MISSING_BACKEND |
| Settings / geral | hostname e idioma | `GET /settings`, `PUT /settings/system` | objeto JSON | campos suportados editáveis; tema permanece local | ADAPT |
| Settings / sistema | versão e uptime | `GET /system/summary` | `SystemSummary` | resumo somente leitura | ADAPT |
| Settings / rede | acesso remoto | `GET /settings`, `PUT /settings/network` | objeto JSON | `remoteAccess` quando presente | ADAPT |
| Settings / updates | canal, automação e verificação | — | — | controles locais temporários | MISSING_BACKEND |
| Settings / manutenção | reboot e shutdown | `POST /system/reboot`, `POST /system/shutdown` | `Job` | confirmação existente + resultado do job | MATCH |
| Settings / diagnóstico | relatório detalhado | — | — | UI demonstrativa preservada | MISSING_BACKEND |
| Files | rota placeholder | — | — | independente; legacy não reintroduzido | MISSING_BACKEND |
| Storage | domínio removido | — | — | nenhum frontend ou endpoint será recuperado | OBSOLETE |

## Decisões

- A UI chama apenas `/api/v1`; o Vite encaminha `/api` para `localhost:18080` em desenvolvimento Windows/WSL.
- DTOs do backend são adaptados para modelos de apresentação somente quando há transformação real.
- `MockProvider` permanece a única implementação exercitada. Linux, Docker, systemd, discos, Samba e host real ficam fora desta fase.
- Metadata editorial da App Store e preferências puramente locais podem permanecer em `data.ts`; estado de sistema, apps e settings não.
