# Installation Guide

## Download Prebuilt Binary (Recommended)

Download from [GitFace Releases](https://github.com/Fish-Ring/GitFace/releases/latest):

| Platform | File |
|---|---|
| Windows amd64 | `gitf-windows-amd64.exe` |
| Windows arm64 | `gitf-windows-arm64.exe` |
| Linux amd64 | `gitf-linux-amd64` |
| Linux arm64 | `gitf-linux-arm64` |
| macOS amd64 | `gitf-darwin-amd64` |
| macOS arm64 | `gitf-darwin-arm64` |

Place the binary in a directory included in your `PATH`.

---

## Build from Source

### Prerequisites

- [Go](https://go.dev/dl/) 1.21 or later

### Windows

```powershell
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

Place `gitf.exe` in `C:\Tools\gitf\` and add to `PATH`:

```powershell
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";C:\Tools\gitf",
    "User"
)
```

### Linux

```bash
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
CGO_ENABLED=0 go build -o gitf .

sudo cp gitf /usr/local/bin/
```

### macOS

```bash
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
CGO_ENABLED=0 go build -o gitf .

sudo cp gitf /usr/local/bin/
```

---

## Verify

```bash
gitf help
```

---

## Update

- **Prebuilt**: Download from [Releases](https://github.com/Fish-Ring/GitFace/releases).
- **Source**: `git pull && go build -o gitf .`

---

> Home: [README](README.md)  
> Usage: [Usage Guide](usage.md)  
> Config: [Config Reference](config.md)
