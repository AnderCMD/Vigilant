# Branching, protection and releases

## Branches

| Branch | Purpose                                                               |
| ------ | --------------------------------------------------------------------- |
| `main` | Always releasable. Protected: only changed through pull requests.     |
| `dev`  | Integration branch. Feature/fix branches are merged here via PR.      |

```text
feature/xyz --PR--> dev --PR (release)--> main --> GitHub Release
```

## Automatic releases

Every push to `main` (i.e. every merged `dev` -> `main` pull request) triggers
`.github/workflows/release.yml`, which:

1. runs the tests,
2. computes the next version from Conventional Commit messages since the last
   tag (`feat:` -> minor, `type!:` or `BREAKING CHANGE` -> major, anything else
   -> patch; the first release is `1.0.0`),
3. builds `vigilant-windows-amd64.exe` and `vigilant-windows-arm64.exe` with the
   version embedded,
4. publishes a GitHub Release with the exes, zip archives and `checksums.txt`,
   plus auto-generated notes.

Use **merge commits** (not squash) for `dev` -> `main` so every commit message
is visible to the version calculation. Update `CHANGELOG.md` before releasing.

## Branch protection (one-time setup)

Protection lives in GitHub settings, not in the repository. After installing and
authenticating the [GitHub CLI](https://cli.github.com) (`gh auth login`), run:

```powershell
./scripts/protect-branches.ps1
```

It creates `dev` if missing and configures:

- **main**: pull request required, `Build & Test` must pass and be up to date,
  conversations resolved, no force-push, no deletion, applies to admins.
- **dev**: `Build & Test` must pass, no force-push, no deletion.
- Repository merge strategy: merge commits only.

> The status check name only appears in GitHub after the CI workflow has run
> once, so push `dev` and let CI finish before running the script.
