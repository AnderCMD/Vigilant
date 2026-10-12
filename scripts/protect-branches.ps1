<#
.SYNOPSIS
  Creates the dev branch if needed and applies branch protection to main and dev.
.DESCRIPTION
  Requires the GitHub CLI (https://cli.github.com) authenticated with admin
  rights on the repository: `gh auth login`.
  - main: pull request required (0 approvals so a solo maintainer can merge),
    "Build & Test" must pass, branch must be up to date, linear-free merge
    commits allowed, no force-push, no deletion, rules apply to admins too.
  - dev:  "Build & Test" must pass, no force-push, no deletion.
#>
param([string]$Repo = "AnderCMD/Vigilant")
$ErrorActionPreference = "Stop"

function Protect($branch, $requirePr) {
  $pr = if ($requirePr) { '{"required_approving_review_count":0,"dismiss_stale_reviews":true}' } else { 'null' }
  $body = @"
{
  "required_status_checks": { "strict": true, "contexts": ["Build & Test"] },
  "enforce_admins": true,
  "required_pull_request_reviews": $pr,
  "restrictions": null,
  "required_conversation_resolution": true,
  "allow_force_pushes": false,
  "allow_deletions": false
}
"@
  $body | gh api --method PUT "repos/$Repo/branches/$branch/protection" --input - | Out-Null
  Write-Host "Protected $branch"
}

$exists = gh api "repos/$Repo/branches/dev" 2>$null
if (-not $exists) {
  $sha = gh api "repos/$Repo/git/ref/heads/main" --jq .object.sha
  gh api --method POST "repos/$Repo/git/refs" -f ref=refs/heads/dev -f sha=$sha | Out-Null
  Write-Host "Created dev from main"
}

Protect "main" $true
Protect "dev" $false

# Merge commits keep Conventional Commit messages visible to the release script.
gh api --method PATCH "repos/$Repo" -F allow_merge_commit=true -F allow_squash_merge=false -F allow_rebase_merge=false -F delete_branch_on_merge=false | Out-Null
Write-Host "Done."
