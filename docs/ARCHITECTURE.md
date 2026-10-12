# Architecture

Vigilant is a small Go program (no CGO) made of a few focused packages.

```text
cmd/vigilant          Entry point: tray menu, wiring, version string
internal/vigil        Platform-neutral scheduler (Runner) that calls a Pulse func
internal/config       JSON settings load/save with validation and defaults
internal/winapi       Win32 calls: SendInput, SetThreadExecutionState, ShellExecute, mutex
internal/autostart    "Start with Windows" via HKCU\...\Run
assets                Embedded tray icons (generated) and the SVG logo
tools/genicon         Renders logo.svg's geometry into .ico/.png files
```

## How a pulse works

1. `vigil.Runner` waits for the configured interval (default 3 minutes).
2. If not paused, it calls `pulse()` in `cmd/vigilant`, which:
   - injects a press/release of **F15** with `SendInput`. Windows counts it as
     user input, so the idle timer resets and chat apps keep you "Active".
     F15 does not exist on common keyboards, so no application reacts to it;
   - calls `SetThreadExecutionState(ES_SYSTEM_REQUIRED [| ES_DISPLAY_REQUIRED])`
     once, which resets the system (and optionally display) idle timers
     without leaving any persistent state behind.
3. Changing the interval or pause state wakes the runner and restarts the countdown.

## Design decisions

- **Go instead of Python**: one static ~3 MB exe, instant start, ~10 MB RAM, no
  interpreter or antivirus-prone PyInstaller bootstrap.
- **Go instead of Rust**: both would work; Go's toolchain is simpler for
  contributors, cross-compiles to arm64 trivially, and the Win32 surface used
  here is tiny.
- **`fyne.io/systray`** talks to the Win32 shell directly (pure Go on Windows).
- **No network, no telemetry**, no key logging: Vigilant only *sends* one key.
- Settings and the autostart entry are per-user; no admin rights are needed.

## Regenerating icons

```powershell
go run ./tools/genicon
```
