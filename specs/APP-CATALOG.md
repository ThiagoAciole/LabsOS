# Catálogo de Apps

O frontend consome `GET /api/v1/catalog/apps`; não acessa CasaOS/ZimaOS diretamente.

Fluxo atual:

`fonte -> downloader -> parser JSON -> validação de id/name -> normalização para catalog.App -> Labs API -> App Store`

No Linux, a fonte padrão é o `gh-pages/index.json` publicado pelo repositório oficial IceWhaleTech/CasaOS-AppStore. O parser aceita o envelope v2 (`apps[]`, `base_url`) e normaliza `title`, `tagline`, `categories`, `icon` e `version` para o contrato LabsOS. O cache em arquivo usa gravação atômica, permissões `0600` e serve como fallback quando a fonte remota estiver indisponível. A fonte externa não é exposta como contrato público do LabsOS.

O catálogo lista metadados de produto. Compose, capabilities, mounts, privileged, Docker socket e outros detalhes de execução não são expostos pela listagem. Instalação real exige um manifest separado e validação específica antes de qualquer operação Docker.
