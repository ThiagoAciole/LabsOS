# Plano de implementação: migração do backend LabsOS

> **Para agentes executores:** SUB-SKILL OBRIGATÓRIO: usar superpowers:subagent-driven-development (recomendado) ou superpowers:executing-plans para implementar este plano tarefa por tarefa. As etapas usam caixas de seleção (`- [ ]`) para acompanhamento.

**Objetivo:** Consolidar o frontend atual na raiz oficial e migrar uma Labs API Go autocontida, executável em modo mock, sem integrar o React.

**Arquitetura:** O servidor HTTP depende de um contrato `Provider`; a inicialização escolhe Mock ou Linux por `LABSOS_MODE`. O mock mantém estado em memória e o Linux permanece explícito como indisponível até existir `labsd`, sem executar comandos do host.

**Stack:** Go 1.26, biblioteca padrão `net/http`, React/TypeScript/Vite preservado.

## Restrições globais

- Frontend em `C:\Projetos\LabsOS\src`.
- Backend em `C:\Projetos\LabsOS\backend`.
- Specs em `C:\Projetos\LabsOS\specs\backend`.
- Nenhuma integração Frontend -> Backend nesta etapa.
- Nenhum shell genérico, Docker real, reboot, shutdown ou acesso destrutivo.
- Apps são produtos; detalhes de containers não fazem parte do contrato público.

---

### Tarefa 1: consolidar a raiz oficial

**Arquivos:** mover o conteúdo de `C:\Projetos\LabsOS\labs-os-panel` para `C:\Projetos\LabsOS`, preservando `.git`, arquivos rastreados, não rastreados e modificações locais.

- [ ] Verificar que a raiz contém somente `labs-os-panel` e `docs` criado por este plano.
- [ ] Mover os itens sem sobrescrever destinos existentes.
- [ ] Confirmar `git status --short` e a presença de `src`, `public`, `package.json` e `vite.config.ts` na raiz.

### Tarefa 2: criar o contrato e o mock por TDD

**Arquivos:**
- Criar: `backend/internal/platform/provider.go`
- Criar: `backend/providers/mock/provider.go`
- Criar: `backend/providers/linux/provider.go`
- Testar: `backend/providers/mock/provider_test.go`

**Interfaces:**
- Produz: `platform.Provider` com System, Apps, AppAction, Settings, UpdateSettings, Events e Job.

- [ ] Escrever testes de estado de Apps, cópia defensiva de Settings e seleção de provider.
- [ ] Executar `go test ./providers/...` e confirmar falha por implementação ausente.
- [ ] Implementar o mínimo para passar, sem dependências externas.
- [ ] Executar `go test ./providers/...` e confirmar sucesso.

### Tarefa 3: expor a API HTTP por TDD

**Arquivos:**
- Criar: `backend/internal/api/server.go`
- Criar: `backend/internal/api/server_test.go`
- Criar: `backend/cmd/labs-api/main.go`
- Criar: `backend/go.mod`

**Interfaces:**
- Consome: `platform.Provider`.
- Produz: `api.New(provider) http.Handler`.

- [ ] Escrever testes de rotas, métodos, JSON, 404 e ações de energia mock seguras.
- [ ] Executar `go test ./internal/api` e confirmar falha por implementação ausente.
- [ ] Implementar handlers e erros estruturados mínimos.
- [ ] Executar `go test ./internal/api` e confirmar sucesso.

### Tarefa 4: limpar e documentar

**Arquivos:** criar os sete documentos em `specs/backend` e atualizar o `README.md` raiz sem alterar o frontend.

- [ ] Executar `gofmt` no backend.
- [ ] Procurar paths antigos, shell genérico e detalhes Docker públicos.
- [ ] Documentar somente endpoints realmente implementados com status READY/PARTIAL/PLANNED/LEGACY.
- [ ] Registrar classificação KEEP/ADAPT/REWRITE/DROP e pendências no relatório.

### Tarefa 5: validar a entrega

- [ ] Executar `go test ./...` em `backend`.
- [ ] Executar `go build ./...` em `backend`.
- [ ] Iniciar com `LABSOS_MODE=mock`, consultar health/system/apps/settings/events e encerrar o processo.
- [ ] Confirmar que o frontend permanece sem alterações causadas pela migração.
