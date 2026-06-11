# 설치 가이드

## 프리빌드 바이너리 다운로드 (권장)

[GitFace Releases](https://github.com/Fish-Ring/GitFace/releases/latest)에서 다운로드:

| 플랫폼 | 파일 |
|---|---|
| Windows amd64 | `gitf-windows-amd64.exe` |
| Windows arm64 | `gitf-windows-arm64.exe` |
| Linux amd64 | `gitf-linux-amd64` |
| Linux arm64 | `gitf-linux-arm64` |
| macOS amd64 | `gitf-darwin-amd64` |
| macOS arm64 | `gitf-darwin-arm64` |

`PATH`에 포함된 디렉토리에 배치하세요.

---

## 소스에서 빌드

### 사전 요구사항

- [Go](https://go.dev/dl/) 1.21 이상

### Windows

```powershell
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

`gitf.exe`를 `C:\Tools\gitf\`에 배치하고 `PATH`에 추가:

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

## 확인

```bash
gitf help
```

---

## 업데이트

- **프리빌드**: [Releases](https://github.com/Fish-Ring/GitFace/releases)에서 다운로드.
- **소스**: `git pull && go build -o gitf .`

---

> 홈: [README](README.md)
> 사용법: [Usage Guide](usage.md)
> 설정: [Config Reference](config.md)
