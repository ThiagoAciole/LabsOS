# Plano prioritário: instalação persistente e updates sem reinstalar a ISO

**Status:** prioritário

## Objetivo

Entregar primeiro uma base de appliance que possa:

1. iniciar pela ISO;
2. instalar o LabsOS uma única vez em um disco de teste;
3. criar e preservar `/DATA` e `/var/lib/labsos`;
4. executar first boot de forma idempotente;
5. iniciar API, dashboard, `labsd`, Docker e kiosk;
6. receber novas versões dos componentes;
7. atualizar com backup, health check e rollback;
8. continuar usando os mesmos dados sem reinstalar a ISO.

A ISO deve ser o mecanismo de instalação inicial, não o mecanismo normal de
atualização do sistema.

## Resultado esperado

```text
ISO de instalação
    ↓ uma única vez
disco persistente
    ├── sistema base
    ├── /DATA
    ├── /var/lib/labsos
    └── /opt/labsos/releases
             ↓
      update versionado
             ↓
      nova release + rollback
```

## Decisões fundamentais

- `/DATA` nunca é apagado durante update.
- `/var/lib/labsos` contém estado, banco SQLite, jobs, eventos e metadados.
- `/opt/labsos/releases/<version>` contém releases imutáveis.
- `/opt/labsos/current` aponta para a release ativa.
- A ISO instala o sistema base e os serviços iniciais.
- Updates substituem a release ativa, não o disco inteiro.
- O update deve funcionar localmente, via arquivo assinado ou via registry.
- Docker volumes persistentes ficam fora da release.
- Qualquer operação destrutiva exige dispositivo explicitamente selecionado.

## Fase 1 — fechar o layout persistente (meio dia)

Definir e documentar:

```text
/boot
/
/DATA
/var/lib/labsos/
├── labsos.db
├── install-state
├── secrets/
├── apps/
├── backups/
└── events/
/opt/labsos/
├── releases/
├── current -> releases/<version>
└── runtime/
```

O layout deve ser usado igualmente no Debian instalado e no ambiente QEMU.

### Critério de aceite

- first boot não recria nem limpa dados existentes;
- restart do serviço preserva SQLite, secrets e `/DATA`;
- um update de teste não altera dados persistentes.

## Fase 2 — installer real e seguro (1–2 dias)

Consolidar o installer em um contrato próprio:

```text
installer.status
installer.disks
installer.validate
installer.start
installer.job
installer.cancel
installer.reboot
```

Fluxo:

```text
listar discos
→ usuário seleciona explicitamente
→ validar proteção
→ particionar somente o alvo
→ criar filesystem
→ instalar rootfs/pacotes
→ criar /DATA
→ criar /var/lib/labsos
→ gravar install-state
→ habilitar first boot
→ reboot
```

### Critério de aceite

- instalação em qcow2 novo;
- nenhum write em disco não selecionado;
- cancelamento deixa estado recuperável;
- falha é registrada no job;
- reboot volta pelo disco instalado.

## Fase 3 — first boot idempotente (meio dia)

O first boot deve:

- validar o install-state;
- criar diretórios persistentes;
- gerar secrets somente se ausentes;
- configurar hostname;
- inicializar SQLite;
- habilitar Docker;
- habilitar `labs-api`, `labsd`, dashboard e kiosk;
- marcar `firstboot.done` somente no final;
- registrar logs no journal.

Executar duas vezes deve produzir o mesmo resultado sem apagar dados.

## Fase 4 — sistema de releases (1 dia)

Criar um release manager simples:

```text
release prepare <version> <artifact>
release verify <version>
release activate <version>
release current
release rollback
```

Cada release deve conter:

```text
release.json
bin/
web/
manifests/
checksums.txt
```

O `current` deve ser atualizado atomicamente por symlink ou mecanismo
equivalente. A versão anterior deve permanecer disponível até a nova release
passar pelo health check.

## Fase 5 — update seguro (1–2 dias)

Fluxo mínimo:

```text
baixar release
→ verificar checksum/assinatura
→ verificar compatibilidade
→ criar backup do estado
→ instalar release ao lado da atual
→ executar migrations
→ ativar serviços
→ health check
→ trocar current
→ registrar sucesso
```

Em caso de falha:

```text
health check falha
→ restaurar current anterior
→ reiniciar serviços
→ marcar update como failed
→ preservar logs para diagnóstico
```

Não atualizar `/DATA` automaticamente nessa fase.

## Fase 6 — SQLite e estado persistente (1 dia)

Depois que o layout estiver estável, colocar no SQLite:

- versão instalada;
- releases disponíveis;
- jobs de instalação/update;
- serviços instalados;
- estado do installer;
- eventos;
- notificações;
- migrations executadas.

O banco deve ter migrations próprias e backup antes de update.

## Fase 7 — Docker e serviços básicos (1–2 dias)

Com o sistema atualizável funcionando, implementar o runtime piloto:

```text
ServicePackage
→ DockerRuntime
→ Docker Compose SDK
→ Docker Engine
```

Começar com Jellyfin ou Syncthing. O objetivo é provar que o update do LabsOS
e o update do serviço são ciclos diferentes:

```text
update do sistema ≠ update do app
```

## Fase 8 — backups e pacote distribuível (depois da base)

Adicionar Restic para backup por serviço e por estado do sistema.

Depois adicionar:

- OCI/ORAS para releases e pacotes;
- Cosign para assinatura;
- registry local;
- sideload assinado;
- registry remoto.

Não bloquear a instalação inicial esperando o registry público.

## Ordem de execução recomendada

```text
1. layout persistente
2. installer em qcow2
3. first boot idempotente
4. releases versionadas
5. update + rollback
6. SQLite persistente
7. DockerRuntime piloto
8. Restic
9. OCI/ORAS
10. Cosign/registry
11. ServicePackage completo
12. expansão para outros serviços
```

## Definition of Done da fundação

Uma pessoa deve conseguir:

```text
make iso-dev
→ iniciar ISO no QEMU
→ instalar em qcow2
→ reiniciar
→ usar o dashboard
→ criar dados em /DATA
→ executar um update da release
→ reiniciar novamente
→ confirmar que os dados continuam presentes
→ provocar falha de health check
→ executar rollback
```

Sem essa sequência funcionando, não considerar a base do appliance pronta.
