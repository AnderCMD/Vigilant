# Configuration

Settings are stored in `%APPDATA%\Vigilant\config.json` and are normally
changed from the tray menu. You can also edit the file while Vigilant is closed
(tray menu -> *Open settings folder*).

```json
{
  "intervalSeconds": 180,
  "keepDisplayAwake": true,
  "startPaused": false
}
```

| Key                | Default | Description                                                           |
| ------------------ | ------- | --------------------------------------------------------------------- |
| `intervalSeconds`  | `180`   | Seconds between pulses. Clamped to 10 - 3600.                         |
| `keepDisplayAwake` | `true`  | Also keep the screen on, not just the system.                         |
| `startPaused`      | `false` | Start in the paused state.                                            |

Invalid or corrupt files are ignored and the defaults are used.

## Start with Windows

The tray menu entry *Start with Windows* adds or removes the value `Vigilant`
under `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`. If you move the
executable, toggle the option off and on again to update the path.
Uninstalling is just deleting the exe (disable autostart first).
