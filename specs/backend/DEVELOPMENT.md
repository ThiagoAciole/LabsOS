# Desenvolvimento do backend

Requisito validado: Go 1.26.

## Executar no Windows

```powershell
cd C:\Projetos\LabsOS\backend
go run ./cmd/labs-api
```

Se a porta 8080 estiver reservada no Windows:

```powershell
$env:LABSOS_ADDR = "127.0.0.1:18080"
go run ./cmd/labs-api
```

Saída esperada:

```text
Labs API http://localhost:8080
```

## Frontend integrado

Com a API em `127.0.0.1:18080`, em outro terminal:

```powershell
cd C:\Projetos\LabsOS
pnpm dev
```

Abra `http://localhost:5173`. O frontend chama `/api/v1` e o proxy Vite encaminha `/api` para a porta 18080.

Em outro terminal:

```powershell
Invoke-RestMethod http://localhost:8080/api/v1/system/health
```

## Validar

```powershell
cd C:\Projetos\LabsOS\backend
gofmt -w cmd internal providers
go test ./...
go build ./...
```

Na raiz, valide o frontend com `pnpm typecheck`, `pnpm lint` e `pnpm build`.

## LinuxProvider no WSL2

O desenvolvimento atual usa Ubuntu WSL2 com systemd. O Go do Windows compila o binário Linux, que é instalado como serviço:

```powershell
cd C:\Projetos\LabsOS\backend
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o ..\.tmp-runtime\labs-api-linux ./cmd/labs-api
wsl -d Ubuntu -u root -- install -m 0755 /mnt/c/Projetos/LabsOS/.tmp-runtime/labs-api-linux /usr/local/bin/labs-api
wsl -d Ubuntu -u root -- install -m 0644 /mnt/c/Projetos/LabsOS/backend/deploy/labs-api-wsl.service /etc/systemd/system/labs-api.service
wsl -d Ubuntu -u root -- systemctl daemon-reload
wsl -d Ubuntu -u root -- systemctl enable --now labs-api.service
```

Validar com `Invoke-RestMethod http://localhost:18080/api/v1/system/summary`. O proxy Vite usa `http://localhost:18080` para funcionar tanto com a API Windows quanto com o encaminhamento do WSL.

## Debian real via SSH

O servidor `labsos` executa `/usr/local/bin/labs-api` como usuário não privilegiado `agent`. A unit versionada é `backend/deploy/labs-api.service`, escuta somente em `127.0.0.1:8080` por padrão e protege operações sensíveis.

No Windows, abra o túnel antes do frontend:

```powershell
ssh -i "$env:USERPROFILE\.ssh\labsos_agent" -N -L 127.0.0.1:18080:127.0.0.1:8080 agent@192.168.0.2
```

O frontend continua usando `/api/v1` e o proxy Vite continua apontando para `localhost:18080`. A API não é publicada diretamente na LAN.

O backend é autocontido nessa pasta; o `labsd` fornece as operações de apps Docker quando instalado.
