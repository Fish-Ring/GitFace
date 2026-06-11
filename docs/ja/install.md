# インストールガイド

## プリビルドバイナリのダウンロード（推奨）

[GitFace Releases](https://github.com/Fish-Ring/GitFace/releases/latest) からダウンロード：

| プラットフォーム | ファイル |
|---|---|
| Windows amd64 | `gitf-windows-amd64.exe` |
| Windows arm64 | `gitf-windows-arm64.exe` |
| Linux amd64 | `gitf-linux-amd64` |
| Linux arm64 | `gitf-linux-arm64` |
| macOS amd64 | `gitf-darwin-amd64` |
| macOS arm64 | `gitf-darwin-arm64` |

`PATH` に含むディレクトリに配置してください。

---

## ソースからビルド

### 前提条件

- [Go](https://go.dev/dl/) 1.21 以降

### Windows

```powershell
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

`gitf.exe` を `C:\Tools\gitf\` に配置し、`PATH` に追加：

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

## 確認

```bash
gitf help
```

---

## 更新

- **プリビルド**: [Releases](https://github.com/Fish-Ring/GitFace/releases) からダウンロード。
- **ソース**: `git pull && go build -o gitf .`

---

> ホーム: [README](README.md)
> 使い方: [Usage Guide](usage.md)
> 設定: [Config Reference](config.md)
