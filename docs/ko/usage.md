# 사용법 가이드

## CLI 명령어

| 명령어 | 설명 |
|---|---|
| `gitf` | TUI 실행 |
| `gitf tui` | TUI 실행 (명시적) |
| `gitf switch <id>` | ID로 프로필 전환 |
| `gitf status` | 현재 신분과 원격 표시 |
| `gitf edit` | 에디터로 설정 열기 |
| `gitf help` | 도움말 표시 |

### CLI 예시

```bash
gitf status
gitf switch personal
gitf edit
```

---

## TUI 인터페이스

`gitf` (인수 없음)로 TUI 실행.

### 레이아웃

```
┌──────────────────────────────────────────────┐
│ GitFace v1.0  │  main [5]                    │
├──────────────────────────────────────────────┤

  【상태】
  • 브랜치:  main (2 파일 언스테이지)
  • 원격:  git@github.com:user/repo.git
  • 신분:  iol <iol@personal.com> [Personal]

  ──────────────────────────────────────────────

  【작업】
  [1] Personal (GitHub)
  [2] Company (GitLab)
  [C] 커밋
  [A] 계정 관리
  [P] 프로바이더 관리
  [S] 설정

  [Y]복사  [R]새로고침  [Q]종료
```

### 핫키 (메인 화면)

| 키 | 동작 |
|---|---|
| `[1]` `[2]` ... | 번호로 신분 전환 |
| `[C]` | 규범 커밋 |
| `[A]` | 계정 관리 |
| `[P]` | 프로바이더 관리 |
| `[S]` | 설정 (언어 전환 + 설정 편집) |
| `[Y]` | 신분 정보를 클립보드에 복사 |
| `[R]` | 상태 새로고침 |
| `[Q]` | 종료 |
| ↑ ↓ / 스크롤 | 탐색 |

### 계정 관리

`[A]`로 계정 관리:

| 키 | 동작 |
|---|---|
| `[N]` | 새 계정 |
| `[D]` | 삭제 모드 ([Enter]로 확인) |
| `[Enter]` | 선택한 계정 편집 |
| `[Esc]` | 뒤로 |
| ↑ ↓ / 스크롤 | 탐색 |

### 계정 양식 (6개 필드)

| 필드 | 설명 |
|---|---|
| ID | `switch <id>`용 고유 ID |
| Name | 표시 이름 |
| Git Name | `git config user.name` |
| Git Email | `git config user.email` |
| Provider | 호스팅 프로바이더 — 입력 또는 `[Ctrl+P]`로 선택 |
| SSH Key | 개인키 경로 — 입력 또는 `[Ctrl+O]`로 스캔 |

### 프로바이더 관리

`[P]`로 프로바이더 관리:

| 키 | 동작 |
|---|---|
| `[N]` | 새 프로바이더 |
| `[D]` | 삭제 모드 |
| `[Enter]` | 선택한 프로바이더 편집 |
| `[Esc]` | 뒤로 |

기본 프로바이더: GitHub / Gitee / GitLab / Bitbucket.

### 설정

`[S]`로 설정:

- **언어** — ← → / Tab으로 전환, Enter로 확인
- **설정 편집** — 에디터로 설정 파일 열기 (저장 후 자동 리로드)

### SSH 키 스캔

1. 계정 양식의 SSH Key 필드에 포커스
2. `[Ctrl+O]` 누르기
3. ↑↓ 또는 스크롤로 키 선택
4. `[Enter]`로 경로 입력

### 프로바이더 선택기

1. 계정 양식의 Provider 필드에 포커스
2. `[Ctrl+P]` 누르기
3. ↑↓ 또는 스크롤로 프로바이더 선택
4. `[Enter]`로 이름 입력

---

## 규범 커밋 흐름

1. `[C]`로 커밋 유형 선택
2. 유형 선택：
   - `f` → `feat` (새 기능)
   - `x` → `fix` (버그 수정)
   - `d` → `docs` (문서)
   - `r` → `refactor` (리팩토링)
3. 설명 입력, `Enter`로 확인
4. 자동 `git add .` → `git commit -m "type: desc"` → `git push`
5. 결과 확인, 아무 키나 눌러 돌아가기

---

## 다중 언어

`LANG` 환경변수에서 자동 감지：

```bash
export LANG=ko_KR.UTF-8   # 한국어
export LANG=en_US.UTF-8   # English
export LANG=zh_CN.UTF-8   # 中文
export LANG=ja_JP.UTF-8   # 日本語
```

설정에서 강제 지정：

```json
{ "lang": "ko" }      // 한국어
{ "lang": "en" }      // English
{ "lang": "zh-CN" }   // 中文
{ "lang": "ja" }      // 日本語
```

우선순위: 설정파일 > `LANG` > `LC_ALL` > `LANGUAGE` > English

---

> 홈: [README](README.md)
> 설치: [Install Guide](install.md)
> 설정: [Config Reference](config.md)
