> [简体中文](../../README.md) · [English](../en/README.md) · [繁體中文](../zh-TW/README.md) · [日本語](../ja/README.md)

# GitFace

Git 다중 ID 및 원격 관리자 — ID, 호스팅 제공자, SSH 키를 한 번에 전환. TUI + CLI 듀얼 인터페이스.

[최신 버전 다운로드](https://github.com/Fish-Ring/GitFace/releases/latest)

---

## 기능

- **프로필 전환** — `user.name`/`user.email` + 원격 저장소 URL + SSH 키를 한 번에 전환
- **호스팅 제공자** — GitHub / Gitee / GitLab / Bitbucket 내장, 자유롭게 추가 가능
- **SSH 키 자동 주입** — 전환 시 `core.sshCommand` 자동 설정, `~/.ssh/config` 불필요
- **SSH 키 스캔** — `[Ctrl+O]`로 `~/.ssh/` 디렉토리를 스캔하고 키 경로를 즉시 입력
- **저장소 설정** — 저장소별 원격 경로 매핑, TUI에서 `[G]`로 편집
- **정형화된 커밋** — 대화형 Conventional Commit (feat/fix/docs/refactor), 자동 add/commit/push
- **릴리스 태그** — TUI에서 `[T]`를 눌러 릴리스 태그를 생성하고 푸시하거나 CLI로 실행
- **브랜치 전환** — TUI에서 `[B]`를 눌러 로컬 브랜치 전환
- **듀얼 인터페이스** — TUI 대화형 인터페이스 + CLI 서브커맨드
- **마우스 스크롤** — 모든 목록에서 스크롤 휠 내비게이션 지원
- **다국어** — 简体中文 / 繁體中文 / English / 日本語 / 한국어 자동 감지
- **제로 설정** — 첫 실행 시 자동으로 설정 파일 생성

---

## 빠른 설치

```powershell
# 소스에서 빌드 (Go 필요)
git clone https://github.com/Fish-Ring/GitFace.git
cd GitFace
go build -o gitf.exe .
```

> 플랫폼별 상세 설치: [설치 가이드](install.md)

---

## 빠른 시작

```bash
# 임의의 Git 저장소에서 실행
gitf                # TUI 시작
gitf status         # Git 상태 보기
gitf switch personal  # 빠른 프로필 전환
```

---

## 문서

[설치](install.md) · [사용법](usage.md) · [설정](config.md)

---

## 라이선스

MIT
