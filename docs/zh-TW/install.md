# 安裝指南

## 下載預編譯二進位（推薦）

從 [GitFace Releases](https://github.com/Fish-Ring/GitFace/releases/latest) 下載對應平台的二進位檔案：

| 平台 | 檔案 |
|---|---|
| Windows amd64 | `gitf-windows-amd64.exe` |
| Windows arm64 | `gitf-windows-arm64.exe` |
| Linux amd64 | `gitf-linux-amd64` |
| Linux arm64 | `gitf-linux-arm64` |
| macOS amd64 | `gitf-darwin-amd64` |
| macOS arm64 | `gitf-darwin-arm64` |

下載後將其放入 `PATH` 環境變數包含的目錄即可。

---

## 從原始碼建置

### 前置要求

- [Go](https://go.dev/dl/) 1.21 或更高版本

### Windows

```powershell
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

建議將 `gitf.exe` 放到 `C:\Tools\gitf\` 目錄下，然後將該目錄加入 `PATH`：

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

## 驗證

```bash
gitf help
```

---

## 更新

- **預編譯二進位**：從 [Releases](https://github.com/Fish-Ring/GitFace/releases) 下載新版替換。
- **原始碼建置**：`git pull && go build -o gitf .`

---

> 首頁：[README](../../README.md)
> 使用指南：[使用指南](usage.md)
> 配置參考：[配置參考](config.md)
