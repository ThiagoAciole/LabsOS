---
name: LabsOS
description: Painel doméstico silencioso para operar apps, arquivos e saúde do servidor.
colors:
  operational-violet: "oklch(0.62 0.26 292)"
  operational-violet-muted: "oklch(0.22 0.08 292)"
  system-black: "oklch(0.035 0 0)"
  system-surface: "oklch(0.125 0 0)"
  system-popover: "oklch(0.16 0 0)"
  telemetry-white: "oklch(0.97 0 0)"
  telemetry-gray: "oklch(0.62 0 0)"
  system-border: "oklch(1 0 0 / 10%)"
  destructive: "oklch(0.704 0.191 22.216)"
typography:
  headline:
    fontFamily: "Geist Variable, sans-serif"
    fontSize: "1.875rem"
    fontWeight: 500
    lineHeight: 1.2
    letterSpacing: "normal"
  title:
    fontFamily: "Geist Variable, sans-serif"
    fontSize: "1rem"
    fontWeight: 500
    lineHeight: 1.5
    letterSpacing: "normal"
  body:
    fontFamily: "Geist Variable, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 400
    lineHeight: 1.5
    letterSpacing: "normal"
  label:
    fontFamily: "Geist Variable, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 500
    lineHeight: 1.4
    letterSpacing: "normal"
rounded:
  sm: "0.525rem"
  md: "0.7rem"
  lg: "0.875rem"
  xl: "1.225rem"
spacing:
  xs: "0.25rem"
  sm: "0.5rem"
  md: "1rem"
  lg: "1.5rem"
  xl: "2rem"
components:
  button-primary:
    backgroundColor: "{colors.operational-violet}"
    textColor: "{colors.telemetry-white}"
    typography: "{typography.body}"
    rounded: "{rounded.md}"
    padding: "0.5rem 0.75rem"
    height: "2.25rem"
  button-secondary:
    backgroundColor: "{colors.operational-violet-muted}"
    textColor: "{colors.telemetry-white}"
    typography: "{typography.body}"
    rounded: "{rounded.md}"
    padding: "0.5rem 0.75rem"
    height: "2.25rem"
  card:
    backgroundColor: "{colors.system-surface}"
    textColor: "{colors.telemetry-white}"
    rounded: "{rounded.xl}"
    padding: "1.5rem"
  input:
    backgroundColor: "{colors.system-popover}"
    textColor: "{colors.telemetry-white}"
    typography: "{typography.body}"
    rounded: "{rounded.md}"
    padding: "0.25rem 0.75rem"
    height: "2.25rem"
---

# Design System: LabsOS

## Overview

**Creative North Star: "Painel de Controle Silencioso"**

O LabsOS usa uma linguagem sóbria, técnica e acolhedora. A interface comunica domínio do sistema sem expor sua complexidade: superfícies escuras, informação organizada e estados claros fazem o servidor parecer um appliance confiável.

A densidade é moderada e orientada à operação. Cor, elevação e movimento aparecem para esclarecer hierarquia ou resposta, nunca para decorar. O resultado deve permanecer próximo de um produto pronto e distante de ferramentas de infraestrutura brutas.

**Key Characteristics:**
- Fundo quase preto e superfícies tonais discretas.
- Violeta reservado para ação, seleção e estado ativo.
- Tipografia Geist compacta, legível e sem teatralidade.
- Bordas sutis, sombras baixas e controles contidos.

## Colors

A paleta combina o Preto de Sistema com o Cinza de Telemetria e usa o Violeta Operacional como sinal funcional raro.

### Primary
- **Violeta Operacional:** ações principais, foco, progresso e navegação ativa.
- **Violeta Operacional Contido:** fundos selecionados e estados secundários sem competir com o conteúdo.

### Neutral
- **Preto de Sistema:** plano-base contínuo da aplicação.
- **Superfície de Sistema:** cards e regiões operacionais sobre o plano-base.
- **Popover de Sistema:** menus, diálogos e superfícies temporárias.
- **Branco de Telemetria:** texto principal e ícones de alta prioridade.
- **Cinza de Telemetria:** metadados, descrições e valores secundários.
- **Borda de Sistema:** separação de baixa intensidade entre superfícies.

### Named Rules

**The Signal, Not Paint Rule.** O violeta identifica ação ou estado; não cobre grandes áreas nem funciona como decoração.

**The Quiet Status Rule.** Verde, amarelo e vermelho aparecem somente quando carregam significado operacional real.

## Typography

**Display Font:** Geist Variable (com fallback sans-serif)
**Body Font:** Geist Variable (com fallback sans-serif)

**Character:** neutra, precisa e contemporânea. Uma única família reduz ruído e mantém títulos, dados e controles no mesmo vocabulário.

### Hierarchy
- **Headline** (500, 1.875rem, 1.2): títulos principais de páginas.
- **Title** (500, 1rem, 1.5): títulos de cards, grupos e ações importantes.
- **Body** (400, 0.875rem, 1.5): conteúdo operacional e descrições.
- **Label** (500, 0.75rem, 1.4): estados, metadados e legendas compactas.

### Named Rules

**The Operational Scale Rule.** Títulos orientam a varredura, mas nunca dominam a área de trabalho como texto promocional.

## Layout

O shell usa sidebar recolhível de 17.5rem no desktop e gaveta de 18rem no mobile. O conteúdo principal ocupa a largura restante, com margem externa de 1rem no desktop e superfície inset arredondada.

As páginas trabalham com ritmo de 1rem a 2rem, grids que passam de uma para duas ou quatro colunas conforme o espaço e padding de 1rem no mobile e 2rem em telas médias. A informação mais importante aparece primeiro; detalhes secundários ficam em cards ou diálogos, não em novas camadas de navegação.

**The Remote-First Rule.** Toda composição deve continuar operável em desktop, tablet e celular sem esconder ações essenciais em hover.

## Elevation & Depth

O sistema usa camadas tonais como principal mecanismo de profundidade. Bordas translúcidas definem limites; sombras pequenas sustentam cards e o shell inset. Sombras maiores ficam restritas a popovers, diálogos e superfícies temporárias.

### Shadow Vocabulary
- **Superfície baixa** (`box-shadow` padrão pequeno): cards e shell inset.
- **Superfície temporária** (`box-shadow` grande): menus, selects, toasts e diálogos.

**The Layer Before Shadow Rule.** Diferencie primeiro pelo tom e pela borda; eleve somente quando a superfície realmente está acima de outra.

## Shapes

Controles usam cantos suavemente curvos; cards e regiões maiores recebem raio mais generoso. Formas circulares ficam reservadas a indicadores, avatares, switches e controles cujo comportamento exige essa silhueta. Bordas são finas e discretas.

**The Contained Curve Rule.** O raio suaviza superfícies operacionais sem transformá-las em cápsulas decorativas.

## Components

### Buttons
- **Shape:** cantos contidos, altura compacta e ícones de 1rem.
- **Primary:** violeta com texto claro; reservado ao comando principal do contexto.
- **Hover / Focus:** alteração tonal no hover e anel violeta translúcido de 3px no teclado.
- **Secondary / Ghost:** superfícies tonais ou transparentes para comandos recorrentes e de menor prioridade.

### Chips
- **Style:** altura compacta, texto pequeno e raio contido.
- **State:** cor semântica somente quando expressa estado; filtros neutros permanecem tonais.

### Cards / Containers
- **Corner Style:** curva ampla e consistente.
- **Background:** Superfície de Sistema, com borda de baixa intensidade.
- **Shadow Strategy:** elevação baixa; popovers assumem a elevação forte.
- **Internal Padding:** geralmente 1.5rem, reduzido em composições densas.

### Inputs / Fields
- **Style:** fundo tonal translúcido, borda sutil e altura de 2.25rem.
- **Focus:** borda Violeta Operacional e anel translúcido de 3px.
- **Error / Disabled:** vermelho semântico para erro; opacidade reduzida e cursor bloqueado para disabled.

### Navigation
- A sidebar usa ícone e texto, recolhe para ícones no desktop e vira gaveta no mobile. O item ativo recebe fundo violeta contido e texto violeta; hover segue o mesmo vocabulário com menor ênfase.

### Metric Cards
- Métricas combinam título curto, valor legível, unidade discreta e indicador visual mínimo. Dados não ganham cor sem função semântica.

## Do's and Don'ts

### Do:
- **Do** use Violeta Operacional para tornar ação, foco e seleção imediatamente reconhecíveis.
- **Do** mantenha métricas, estados e ações fáceis de varrer em poucos segundos.
- **Do** preserve controles acessíveis por teclado e ações essenciais visíveis em touch.
- **Do** use camadas tonais antes de adicionar sombra ou ornamentação.

### Don't:
- **Don't** reproduza a aparência de painel Linux ou exponha jargão de infraestrutura na interface comum.
- **Don't** transforme o produto em dashboard corporativo denso, com excesso de tabelas, filtros ou KPIs.
- **Don't** use neon, brilhos intensos ou linguagem visual gamer.
- **Don't** use violeta como preenchimento decorativo de grandes superfícies.
