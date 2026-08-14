# Providers

`internal/platform.Provider` define o contrato usado pelos handlers: System, Power, Apps, Settings, Events e Jobs. A seleção ocorre uma vez na inicialização por `LABSOS_MODE`.

## MockProvider

- funcional no Windows sem Debian, Docker, systemd, `/proc`, `/sys` ou hardware;
- mantém apps, settings, jobs e events em memória;
- reboot e shutdown apenas criam jobs com mensagem `simulated in mock mode`;
- não persiste estado após reiniciar o processo.

## LinuxProvider

Implementado somente para leitura de System: hostname, uptime, CPU, RAM, temperatura disponível e rede são obtidos de `/proc`, `/sys` e interfaces do kernel. Apps, Settings, Events e Jobs ainda delegam ao MockProvider para preservar o fluxo de desenvolvimento. Power falha fechado com HTTP 503 e nunca executa reboot ou shutdown.

Esse modo híbrido é temporário e explícito. Docker Apps, discos, Samba, Files e operações privilegiadas continuam fora do LinuxProvider.

## Apps/Docker

Apps são produtos. Nenhum DTO público contém `containerId`, `composeFile`, `dockerNetwork`, `imageSha`, `dockerVolume` ou `containerName`. O futuro provider Linux deve conversar com `labsd`; a API pública não acessará o socket Docker.

## Segurança

- modos desconhecidos impedem a inicialização;
- energia e domínios Linux ainda não implementados falham fechado com HTTP 503;
- não existe endpoint genérico de comandos;
- providers recebem `context.Context`;
- operações privilegiadas futuras exigem métodos tipados e Unix socket local restrito.
