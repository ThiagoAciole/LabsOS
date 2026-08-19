# ISO LabsOS

Builder Debian 13 Trixie baseado em `live-build`, seguindo o modelo de imagem final do StartOS. Funciona em Ubuntu/Debian Linux nativo e também em WSL2; em WSL2, o filesystem Linux costuma ser mais rápido que `/mnt/c`.

```bash
make setup
make packages
make iso
make test
```

`make iso` exige Ubuntu/Debian Linux, internet para o mirror Debian e `live-build`. WSL2 é apenas uma alternativa para quem está no Windows. A imagem é construída diretamente como sistema live/hybrid; ela não usa Debian Installer nem `simple-cdd`. O instalador persistente em disco e o first boot definitivo são etapas seguintes do pipeline LabsOS.

O código pode estar em um disco exFAT/NTFS, mas o diretório temporário do `live-build` precisa estar em um filesystem Linux nativo. Se o clone estiver em `/run/media`, use:

```bash
mkdir -p ~/labsos-image-build
LABSOS_IMAGE_BUILD_DIR="$HOME/labsos-image-build" make iso
```

O build inclui `zstd`, firmware Realtek e o plugin oficial Docker Compose v5.1.4 em `/usr/libexec/docker/cli-plugins/docker-compose`.

Para Windows:

```bash
cp dist/LabsOS-0.1.0-amd64.iso /mnt/c/Users/<USER>/Downloads/
```

Grave a ISO com Rufus/Ventoy. O arquivo qcow2 fica em `build/` e nunca aponta para `/dev/sda` ou `/dev/nvme*`.
