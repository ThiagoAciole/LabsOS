# ISO LabsOS

Builder Debian 13 Trixie para WSL2. O clone em `/mnt/c` funciona, mas o filesystem Linux do WSL é mais rápido.

```bash
make setup
make packages
make iso
make test
```

`make iso` exige WSL2 Debian/Ubuntu, internet para o mirror Debian e `simple-cdd`. A instalação usa Debian Installer, não seleciona disco automaticamente e não inclui apps opcionais. O preseed não contém senhas; a conta `labs` deve ser criada no fluxo do instalador.

Para Windows:

```bash
cp dist/LabsOS-0.1.0-amd64.iso /mnt/c/Users/<USER>/Downloads/
```

Grave a ISO com Rufus/Ventoy. O arquivo qcow2 fica em `build/` e nunca aponta para `/dev/sda` ou `/dev/nvme*`.
