# LabsOS runtime

Área do runtime do appliance LabsOS. O runtime é parte do monólito do produto:
não é um serviço externo nem um projeto independente.

Hoje sua implementação Go permanece em `backend/labsd` e `backend/runtime`
porque esses pacotes compartilham o módulo `backend` e os contratos internos
dos providers. Esta pasta documenta o limite arquitetural para a próxima
extração, quando houver um runtime de containers e um agente privilegiado com
ciclo de vida próprio.

```text
web → backend/internal/api → backend/providers → runtime/labsd
                                           └── runtime/container-runtime
```

Qualquer extração futura deve manter o contrato tipado e o princípio de que o
frontend nunca acessa o runtime diretamente.
