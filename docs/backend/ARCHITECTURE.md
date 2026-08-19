# Arquitetura do backend

```text
Frontend React
  -> HTTP /api/v1
  -> Handler (internal/api)
  -> Provider contract (internal/platform)
  -> LinuxProvider
```

O slice atual não cria uma camada Service: não há regra suficientemente complexa fora do provider. Quando jobs persistentes ou orquestração entre domínios surgirem, Services entram entre handlers e providers.

## Responsabilidades

- **Handler:** valida método, path e JSON; converte resultados e erros para HTTP.
- **Service:** planejado para orquestração comprovadamente necessária; não existe nesta fase.
- **Provider:** contrato estável para operações dependentes do ambiente.
- **Model/DTO:** tipos em `internal/platform`; não expõem IDs, redes ou volumes Docker.
- **Linux:** é o provider operacional padrão; lê dados do host sem privilégios e falha fechado quando uma capacidade não está disponível.
- **Testes:** podem usar providers controlados diretamente, sem seleção de modo no produto ou configuração de runtime.

## Privilégios futuros

```text
labs-api (não-root, LAN)
  -> HTTP/JSON via /run/labsos/labsd.sock
  -> labsd (root, local-only)
  -> operações Linux tipadas
```

`labsd` não está implementado. Não haverá `/shell`, `/exec` ou concatenação de dados do usuário em `sh -c`. O frontend nunca acessará Docker, systemd, `/proc`, `/sys`, discos ou filesystem real diretamente.
