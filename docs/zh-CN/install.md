# 安装指南

## 下载预编译二进制（推荐）

从 [GitFace Releases](https://github.com/Fish-Ring/GitFace/releases/latest) 下载对应平台的二进制文件：

| 平台 | 文件 |
|---|---|
| Windows amd64 | `gitf-windows-amd64.exe` |
| Windows arm64 | `gitf-windows-arm64.exe` |
| Linux amd64 | `gitf-linux-amd64` |
| Linux arm64 | `gitf-linux-arm64` |
| macOS amd64 | `gitf-darwin-amd64` |
| macOS arm64 | `gitf-darwin-arm64` |

下载后将其放入 `PATH` 环境变量包含的目录即可。

---

## 从源码构建

### 前置要求

- [Go](https://go.dev/dl/) 1.21 或更高版本

### Windows

```powershell
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

推荐将 `gitf.exe` 放到 `C:\Tools\gitf\` 目录下，然后将该目录加入 `PATH`：

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

## 验证

```bash
gitf help
```

---

## 更新

- **预编译二进制**：从 [Releases](https://github.com/Fish-Ring/GitFace/releases) 下载新版替换。
- **源码构建**：`git pull && go build -o gitf .`

---

> 首页：[首页](../../README.md)  
> 使用指南：[使用指南](usage.md)  
> 配置参考：[配置参考](config.md)
