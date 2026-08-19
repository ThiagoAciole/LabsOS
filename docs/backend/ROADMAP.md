# Roadmap do backend

1. **Phase 1 — API contratual:** núcleo de System, Apps, Settings, Events e Jobs pronto; fixtures em memória ficam restritas aos testes.
2. **Phase 2 — Frontend integration:** concluída para Home, Apps, App Store e Settings via API local; lacunas estão em `FRONTEND-INTEGRATION.md`.
3. **Phase 3/4 — Linux SystemProvider:** leituras seguras de System, rede e armazenamento vêm do sistema local; capacidades ausentes falham fechado.
4. **Phase 5/6 — Docker Apps e labsd:** IPC HTTP/JSON em Unix socket, privilégios mínimos, jobs persistentes e operações tipadas.
5. **Phase 7–9 — packaging, installer e ISO:** units, pacotes reproduzíveis, installer protegido e validação BIOS/UEFI/QEMU.

Files pode voltar como slice próprio limitado a `/DATA`, com proteção contra traversal e symlink escape. Storage permanece futuro e isolado; nenhum disco secundário será selecionado ou alterado automaticamente.
