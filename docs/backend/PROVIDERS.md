# Providers

`internal/platform.Provider` define o contrato usado pelos handlers: System, Power, Apps, Settings, Events e Jobs. A API inicia com o provider Linux; não há seleção de modo no runtime.

## Fixtures de testes

Os testes unitários podem usar um provider em memória isolado, apenas como fixture.
Ele não é selecionável pelo runtime, não é empacotado como backend operacional e
não altera o comportamento do provider Linux.

## LinuxProvider

O provider Linux lê hostname, uptime, CPU, RAM, temperatura disponível e rede de `/proc`, `/sys` e interfaces do kernel. Apps, Settings, Events e Jobs usam fontes locais; capacidades ausentes falham fechado. Power exige a política de operações reais e nunca executa reboot ou shutdown sem autorização.

Capacidades ausentes falham fechado. Operações potencialmente destrutivas ou
privilegiadas continuam protegidas pela política local e exigem confirmação explícita.

## Apps/Docker

Apps são produtos. Nenhum DTO público contém `containerId`, `composeFile`, `dockerNetwork`, `imageSha`, `dockerVolume` ou `containerName`. O futuro provider Linux deve conversar com `labsd`; a API pública não acessará o socket Docker.

## Segurança

- energia e domínios Linux indisponíveis falham fechado com HTTP 503;
- não existe endpoint genérico de comandos;
- providers recebem `context.Context`;
- operações privilegiadas futuras exigem métodos tipados e Unix socket local restrito.
