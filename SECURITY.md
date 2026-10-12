# Security Policy

## Supported Versions

Vigilant is a small utility distributed through GitHub Releases. Only the
latest release is supported with security updates.

## Reporting a Vulnerability

If you discover a security vulnerability, please report it privately by
opening a [GitHub Security Advisory](../../security/advisories/new) instead
of a public issue. Include:

- A description of the vulnerability and its potential impact
- Steps to reproduce it
- Any relevant logs or environment details

We will do our best to respond promptly and address confirmed issues.

## Scope Notes

Vigilant does not collect data, connect to the internet, or track keystrokes.
It only simulates a single, non-disruptive key press (F15) and resets the
Windows idle timers at a fixed interval. Settings are stored locally in
`%APPDATA%\Vigilant\config.json`; the optional autostart entry lives in the
current user's registry `Run` key. Release binaries ship with SHA-256 checksums. Reports related to this core
behavior being misused (e.g. as a basis for building surveillance tools) are
out of scope, since the project intentionally performs no data collection.
