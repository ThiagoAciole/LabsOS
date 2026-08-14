# Arquitetura do backend

```text
Frontend React
  -> HTTP /api/v1
  -> Handler (internal/api)
  -> Provider contract (internal/platform)
  -> MockProvider | LinuxProvider
```

O slice atual não cria uma camada Service: não há regra suficientemente complexa fora do provider. Quando jobs persistentes ou orquestração entre domínios surgirem, Services entram entre handlers e providers.

## Responsabilidades

- **Handler:** valida método, path e JSON; converte resultados e erros para HTTP.
- **Service:** planejado para orquestração comprovadamente necessária; não existe nesta fase.
- **Provider:** contrato estável para operações dependentes do ambiente.
- **Model/DTO:** tipos em `internal/platform`; não expõem IDs, redes ou volumes Docker.
- **Mock:** simula sistema, apps, settings, events e jobs em memória sem tocar o host.
- **Linux:** lê métricas de System sem privilégios; no WSL, domínios ainda não migrados usam fallback mock explícito e Power falha fechado.

## Privilégios futuros

```text
labs-api (não-root, LAN)
  -> HTTP/JSON via /run/labsos/labsd.sock
  -> labsd (root, local-only)
  -> operações Linux tipadas
```

`labsd` não está implementado. Não haverá `/shell`, `/exec` ou concatenação de dados do usuário em `sh -c`. O frontend nunca acessará Docker, systemd, `/proc`, `/sys`, discos ou filesystem real diretamente.
