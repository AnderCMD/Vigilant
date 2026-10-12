# Contributing to Vigilant

Thanks for your interest in contributing! This guide explains how to get set up and submit changes.

## Branching model

- `main` - stable, protected, every merge publishes a GitHub Release.
- `dev` - integration branch; **all pull requests target `dev`**.

Create your branch from `dev` (`feat/my-feature`, `fix/tray-icon`, ...). Maintainers
periodically open a `dev` -> `main` pull request to cut a release. Details in
[docs/RELEASING.md](docs/RELEASING.md).

## Getting started

1. Install [Go](https://go.dev/dl/) 1.26 or newer (Windows is required to run the app).
2. Fork and clone:
   ```bash
   git clone https://github.com/<your-username>/Vigilant.git
   cd Vigilant
   git checkout dev
   git checkout -b feat/my-feature
   ```
3. Run it: `.\scripts\run.ps1`, build it: `.\scripts\build.ps1`.

## Before opening a pull request

```powershell
gofmt -l .        # must print nothing
go vet ./...
go test ./...
```

CI runs the same checks plus a build for amd64 and arm64.

## Guidelines

- Keep pull requests focused on one change.
- Use [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`,
  `docs:`, `chore:`, `refactor:`, `test:`, `feat!:` for breaking changes). The next
  release version is derived from them.
- Put platform-neutral logic in `internal/vigil` / `internal/config` with unit tests;
  keep Win32 code in `internal/winapi` and `internal/autostart`.
- Update documentation and `CHANGELOG.md` when behavior changes.
- Everything in the repository (code, comments, docs, commits) is written in English.

## Reporting bugs and ideas

Use the [issue templates](.github/ISSUE_TEMPLATE). Include your Windows version,
Vigilant version (shown in the tray menu) and steps to reproduce.

## License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
