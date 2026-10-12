# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project follows
[Semantic Versioning](https://semver.org/). GitHub Releases also contain
auto-generated notes for every version.

## [Unreleased]

### Changed
- Rewritten in Go: a single ~3 MB native executable with no runtime or
  dependencies to install (previously Python + PyInstaller, ~15 MB+).
- Activity is now simulated with the Win32 `SendInput` API (F15 key) and
  `SetThreadExecutionState`, instead of `pyautogui`.
- New logo and tray icons, generated reproducibly by `tools/genicon`.

### Added
- "Start with Windows" option (per-user registry `Run` key).
- Selectable pulse interval (1, 3, 5, 10, 15 minutes) and "Keep display on" option.
- Persistent settings in `%APPDATA%\Vigilant\config.json`.
- Single-instance protection.
- Left-click on the tray icon toggles pause/resume.
- Windows executable metadata (icon, version, description) and arm64 builds.
- CI (format, vet, tests, build), automatic GitHub Releases from `main`,
  Dependabot, CODEOWNERS and branch-protection script.
- Documentation in `docs/`.

### Removed
- Python sources, `requirements.txt` and the `.bat` helper scripts.
