# Shared contracts

Contratos compartilháveis do monólito LabsOS.

Esta área deve receber apenas modelos consumidos por mais de um domínio, como:

- eventos e notificações;
- jobs e progresso;
- manifests e service packages;
- contratos de runtime.

Enquanto houver um único consumidor, o contrato deve permanecer junto do
domínio que o implementa.
