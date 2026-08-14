# System Apps

System Apps são componentes internos mantidos pelo LabsOS. O primeiro contrato é `files`, o File Manager do produto.

User Apps aparecem em `Apps` e na App Store. System Apps não aparecem nessas superfícies e não são removíveis pelo usuário comum.

O lifecycle interno previsto é: `ensure installed`, `start`, `stop`, `restart`, `health` e `url`. A implementação do runtime e do provisioning deve permanecer atrás de interfaces tipadas; não existe endpoint genérico de execução.

O File Manager deverá montar somente `/DATA`, manter sua configuração em área controlada pelo LabsOS e ser acessado pela rota `Files`, sem publicar sua porta diretamente na LAN.

O manifesto de referência está em `specs/files-system-app.compose.yml`: a porta é publicada somente em `127.0.0.1:8081`, o usuário só recebe `/DATA:/srv` e a configuração fica em `/var/lib/labsos/system-apps/files`. O System App está provisionado no Debian e saudável. Em desenvolvimento, o frontend usa `/file-manager/`, encaminhado pelo Vite para `localhost:18081` através do túnel SSH.
