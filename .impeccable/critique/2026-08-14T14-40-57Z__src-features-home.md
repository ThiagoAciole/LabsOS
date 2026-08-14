---
target: src/features/Home
total_score: 21
max_score: 40
na_heuristics: 
p0_count: 1
p1_count: 3
timestamp: 2026-08-14T14-40-57Z
slug: src-features-home
---
# Critique: Home

## Design Health Score

| # | Heuristic | Score | Key Issue |
|---|---|---:|---|
| 1 | Visibility of System Status | 1 | Sem loading, atualização, stale, falha ou progresso de ações. |
| 2 | Match System / Real World | 3 | Linguagem de appliance funciona; alguns rótulos em inglês e telemetria sem contexto. |
| 3 | User Control and Freedom | 2 | Há comandos de apps, mas faltam cancelamento, confirmação e recuperação. |
| 4 | Consistency and Standards | 3 | Componentes coerentes; cores de upload/download conflitam com semântica de status. |
| 5 | Error Prevention | 1 | Parar não parece destrutivo e não há guardas aparentes. |
| 6 | Recognition Rather Than Recall | 3 | Ícones, nomes e ações são reconhecíveis; gráfico de uptime é opaco. |
| 7 | Flexibility and Efficiency | 2 | Abrir app é direto; faltam atalhos e visão compacta. |
| 8 | Aesthetic and Minimalist Design | 3 | Calmo e consistente, mas relógio e cards altos ocupam espaço operacional. |
| 9 | Error Recovery | 1 | Não existem estados de erro, retry ou dados parciais. |
| 10 | Help and Documentation | 2 | Termos comuns são claros; limites e consequências não são explicados. |
| **Total** | | **21/40** | **Acceptable** |

## Design Specificity Verdict

**LLM assessment:** moderadamente específico. Apps, telemetria, uptime e atividade situam o produto no homelab, mas greeting, grid de KPIs e cards uniformes ainda poderiam pertencer a um dashboard genérico. O caráter do LabsOS está mais nos substantivos do que no modelo operacional.

**Deterministic scan:** o CLI retornou zero achados em `src/features/Home`. No browser, a página composta apresentou 15 alertas: 11 `nested-cards`, 2 `layout-transition`, 1 `cramped-padding` e 1 `skipped-heading`. Os 11 nested cards são majoritariamente falsos positivos causados pelo contêiner visual da Home; os 2 transitions pertencem ao shell compartilhado. Heading pulado é real; cramped padding precisa de confirmação no componente.

**Visual overlays:** a injeção funcionou em uma aba isolada `[Human]`; 15 overlays foram renderizados e inspecionados. A aba e os servidores temporários foram encerrados após a coleta.

## Overall Impression

A base visual é coerente, calma e alinhada ao produto. A maior oportunidade não é cosmética: transformar telemetria e controles que parecem confiáveis em estados realmente verificáveis, acionáveis e recuperáveis.

## What's Working

- O conjunto CPU, RAM, disco, rede, uptime, apps e eventos respeita o escopo appliance-first.
- O vocabulário visual é consistente e usa violeta como sinal, não decoração.
- Os cards de app mantêm Abrir visível e escondem comandos secundários sem sobrecarregar.

## Priority Issues

### P0 - Comandos operacionais sem contrato de resultado

**Why it matters:** Abrir, Reiniciar e Parar parecem funcionais, mas não exibem wiring, pending, sucesso ou falha; o fluxo principal pode terminar silenciosamente.

**Fix:** esconder comandos ainda indisponíveis ou conectá-los com pending/disabled, confirmação para Parar e feedback de sucesso/erro.

**Suggested command:** `$impeccable harden src/features/Home`

### P1 - Dados estáticos apresentados como verdade atual

**Why it matters:** métricas, IP, status, data e atividade não distinguem mock, cache, stale, offline ou falha de API.

**Fix:** definir estados `loading`, `ready`, `stale` e `unavailable`; mostrar `Atualizado há...` e nunca mascarar falha com mock em produção.

**Suggested command:** `$impeccable harden src/features/Home`

### P1 - Hierarquia privilegia saudação e relógio

**Why it matters:** durante uma checagem remota, saúde do servidor e apps deveriam responder primeiro; o relógio ocupa o maior peso visual.

**Fix:** elevar um resumo geral de saúde, reduzir Welcome/Clock e compactá-los em telas baixas ou estreitas.

**Suggested command:** `$impeccable layout src/features/Home`

### P1 - Semântica de status fragmentada

**Why it matters:** pontos verdes e setas vermelha/verde dependem de cor e confundem direção de tráfego com erro/sucesso; uptime não possui período nem resumo acessível.

**Fix:** parear cor com texto, neutralizar upload/download e fornecer legenda e resumo textual para histórico.

**Suggested command:** `$impeccable audit src/features/Home`

### P2 - Mobile apenas empilha o desktop

**Why it matters:** cards com altura mínima, relógio grande, timestamps ocultos e botão de 36px tornam a checagem por celular longa e menos informativa.

**Fix:** criar prioridade mobile, compactar apps, manter horários abreviados e elevar alvos touch para 44px.

**Suggested command:** `$impeccable adapt src/features/Home`

## Persona Red Flags

**Alex (power user):** precisa varrer vários cards de mesmo peso para encontrar anomalias; não há visão compacta, atalho evidente ou ação em lote.

**Sam (teclado, leitor de tela e baixa visão):** mudanças de status não são anunciadas, cores carregam significado, o gráfico de uptime não tem resumo útil e o botão de opções tem alvo reduzido.

**Marina (responsável pelo homelab doméstico):** Online não diz o que foi verificado; atividade não oferece evidência; Parar não explica consequência para serviços usados pela casa.

## Minor Observations

- A copy mistura português com `RAM Usage`, `Disk Usage`, `CPU Usage`, `On Time` e `Download/Upload`.
- O DOM pula de `h1` para `h3`; Recent Activity ainda aninha parágrafos dentro de `h3`.
- `<time>` não possui `dateTime`, e horários relativos somem no mobile.
- O detector de browser sinalizou padding reduzido em um botão; a origem provável precisa ser confirmada antes da edição.

## Questions to Consider

1. Qual resposta precisa ser confiável em três segundos: tudo saudável, problemas ativos ou apps disponíveis?
2. O clima merece o primeiro viewport de um appliance de servidor?
3. O que exatamente Online comprova: Agent, LAN, internet ou interface de rede?
4. Qual consequência deve ser explicada antes de parar AdGuard ou Home Assistant?
