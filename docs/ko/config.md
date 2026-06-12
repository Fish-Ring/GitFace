# 설정 레퍼런스

## 설정 파일 위치

- **기본값**: `~/.config/gitface/config.json`
- 첫 실행 시 자동 생성

---

## 구조

```json
{
  "lang": "ko",
  "active_profile_id": "",
  "profiles": [
    {
      "id": "personal",
      "name": "개인 (GitHub)",
      "git_name": "yourname",
      "git_email": "you@personal.com",
      "provider_id": "github",
      "ssh_identity_file": "~/.ssh/id_ed25519_personal"
    }
  ],
  "providers": [
    {
      "id": "github",
      "name": "GitHub",
      "host": "github.com"
    }
  ]
}
```

### 최상위 필드

| 필드 | 타입 | 설명 |
|---|---|---|
| `lang` | string | 언어: `"en"`, `"zh-CN"`, `"ja"`, `"ko"` |
| `active_profile_id` | string | 예약 |
| `profiles` | array | 프로필 목록 |
| `providers` | array | 프로바이더 목록 (기본: GitHub, Gitee, GitLab, Bitbucket) |

### Profile 필드

| 필드 | 설명 |
|---|---|
| `id` | 고유 ID, `switch <id>`에서 사용 |
| `name` | TUI 표시 이름 |
| `git_name` | `git config user.name`에 기록 |
| `git_email` | `git config user.email`에 기록 |
| `provider_id` | 프로바이더의 `id`에 연결; 원격 URL 호스트를 프로바이더의 `host`로 리라이트 |
| `ssh_identity_file` | SSH 개인키 경로; `git config core.sshCommand`로 자동 주입 |

### Provider 필드

| 필드 | 설명 |
|---|---|
| `id` | 고유 ID (예: `github`) |
| `name` | 표시 이름 (예: `GitHub`) |
| `host` | Git 호스트 도메인 (예: `github.com`, `gitee.com`) |

---

## 저장소별 구성

각 저장소는 `.git/gitf.json`에 프로바이더별 원격 경로 매핑을 저장할 수 있습니다. 이 파일은 첫 신분 전환 시 자동 생성되며 git은 무시합니다 (`.git/` 내부).

### 구조

```json
{
  "paths": {
    "github": "Fish-Ring/GitFace",
    "gitee": "Ableand/git-face"
  }
}
```

### 필드

| 필드 | 설명 |
|---|---|
| `paths` | 프로바이더 ID에서 해당 프로바이더의 저장소 경로로의 매핑 |

### 작동 방식

신분 전환 시 GitFace는 원격 URL을 다시 작성합니다. 저장소 구성은 각 프로바이더에 사용할 경로를 지정합니다. 예를 들어 GitHub 저장소가 `Fish-Ring/GitFace`이고 Gitee 미러가 `Ableand/git-face`인 경우, 여기에서 두 경로를 구성합니다.

TUI에서 `[G]`를 눌러 편집하거나, 수동으로 `.git/gitf.json`을 편집하세요.

---

## SSH 설정

### 자동 설정 (권장)
GitFace는 `git config core.sshCommand "ssh -i <키 경로>"`로 SSH 키를 자동 주입합니다. **`~/.ssh/config` 수동 설정이 필요하지 않습니다.** 신분 전환 시 자동 주입되며, 전환 해제 시 자동 클리어됩니다.

### 수동 ~/.ssh/config (선택)
기존 방식으로 `~/.ssh/config`를 설정할 수도 있습니다:

```
Host github-personal
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_personal
```

하지만 GitFace의 프로바이더 시스템은 원격 호스트를 실제 도메인(예: `github.com`)으로 리라이트합니다. 키 선택은 `core.sshCommand`로 자동 처리됩니다.

---

## 신분 전환 시 실행되는 작업

신분으로 전환하면 GitFace가 자동 실행:

1. `git config user.name` ← 프로필의 `git_name`
2. `git config user.email` ← 프로필의 `git_email`
3. 원격 URL 호스트 이름을 프로바이더의 도메인으로 리라이트
4. `git config core.sshCommand "ssh -i <ssh_identity_file>"` ← 키 자동 주입

---

> 홈: [README](README.md)
> 설치: [Install Guide](install.md)
> 사용법: [Usage Guide](usage.md)
