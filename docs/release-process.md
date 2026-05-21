# Release Process

## Overview

hfgo uses an automated release workflow driven by
[release-please](https://github.com/googleapis/release-please)
in **manifest mode** and
[Conventional Commits](https://www.conventionalcommits.org/).

There are two tracks:

- **Mainline releases**: Patch and minor releases from `main` (fully automated)
- **Release candidate (RC) releases**: Major or significant releases from a
  dedicated draft branch (parameterized same workflow)

The workflow uses **separate configuration files** for stable and prerelease
runs, selected programmatically based on the branch name. This avoids the
need for draft branches to recreate or modify configuration.

---

## Configuration Files

All configuration files live on all branches. The workflow selects the
appropriate pair at runtime.

### Stable Config

| File                                    | Purpose                             |
|-----------------------------------------|-------------------------------------|
| `.github/release-please-config.json`    | Stable release-please configuration |
| `.github/.release-please-manifest.json` | Tracks last stable version          |

### Prerelease Config

| File                                    | Purpose                                                                 |
|-----------------------------------------|-------------------------------------------------------------------------|
| `.github/prerelease-please-config.json` | Prerelease configuration (`prerelease: true`, `versioning: prerelease`) |
| `.github/.prerelease-manifest.json`     | Tracks last RC version                                                  |

Both config files use the `packages` key (manifest mode) and include the
same base options (`release-type: go`, `version-file`, `changelog-sections`,
etc.). Only the prerelease-specific flags differ.

### Why Separate Manifests?

Each branch reads its own manifest file. This provides **structural version
isolation**: the draft branch's RC versions never interfere with main's
stable version baseline, and vice versa.

Because both manifest files exist on all branches, a stale stable manifest
(created when the draft branch was cut from `main`) naturally carries the
correct stable baseline. During promotion, no manifest reset is needed.

---

## Mainline Releases

### Trigger

Pushes to `main` trigger the [`Release` workflow](`.github/workflows/release.yml`).
The workflow selects the **stable** config/manifest pair.

### Mainline Workflow

1. Maintainers merge PRs to `main` using squash merge.
2. PR titles follow Conventional Commits (enforced by `pr-title-check.yml`).
3. release-please analyzes commits since the last stable release.
4. A **draft release PR** is created/updated on `main` with:
   - Updated version in `internal/sdkversion/version.go`
   - Updated `.github/.release-please-manifest.json`
   - Auto-generated changelog
5. A maintainer reviews and merges the release PR when ready.
6. release-please creates a GitHub release with the tag and changelog.

### Version Determination

| Bump  | Trigger                                     |
|-------|---------------------------------------------|
| Major | Breaking change (`feat!:` or `fix!:` title) |
| Minor | New feature (`feat:`)                       |
| Patch | Bug fix (`fix:`) or refactoring             |
| None  | Documentation, tests, chore                 |

---

## Release Candidate (RC) Releases

### When to Use

For major version upgrades or significant changes that need testing before a
formal release.

### Branch Convention

- **Major revision**: `vX-draft` (e.g., `v5-draft`)
- **Minor version**: `vX.Y-draft` (e.g., `v5.1-draft`)
- **Patch version**: `vX.Y.Z-draft` (e.g., `v5.1.1-draft`)
  - Note: patches should rarely, if ever, justify tracking a separate
    release-candidate branch

### Setup

1. Create the draft branch from `main`:

   ```bash
   git checkout -b vX-draft main
   ```

1. Update `go.mod` to the new module path (e.g., `github.com/Kardbord/hfgo/vX`).
1. Update `.github/.prerelease-manifest.json` with the starting RC version:

   ```json
   {
     ".": "X.Y.Z-rc.0"
   }
   ```

1. Push the branch.

The stable manifest (`.github/.release-please-manifest.json`) remains at the
last stable version from `main`. It is unused on the draft branch but
preserved for promotion.

### RC Workflow

Pushes to `v*-draft` trigger the same `Release` workflow, which selects the
**prerelease** config/manifest pair.

1. PRs merged to the draft branch accumulate conventionally.
1. release-please bumps the pre-release version (`rc.1`, `rc.2`, ...).
1. Each merge of the RC release PR creates a GitHub **pre-release** with a
   `vX.Y.Z-rc.N` tag.
1. `main` continues to receive patches and minor updates independently.

### Keeping the Draft Branch in Sync

Maintainers can:

- **Cherry-pick** fixes from `main` that the RC should also include.
- **Merge `main` into the draft branch** periodically (this also pulls in fixes
  that should appear in the RC changelog, which is slightly redundant but
  acceptable).

### Promoting an RC to a Formal Release

When the RC is stable:

1. Ensure the final RC PR is merged so the release candidate tag reflects all
   desired changes.
1. Optional: if this is a major revision, deprecate the previous major version
   in `doc.go` and perform a final patch release so that the deprecation is
   reflected by the go module proxy.
1. Create a promotion PR from the draft branch to `main`.
   - Resolve any merge conflicts in favor of the draft branch's code changes.
   - The stale `.github/.release-please-manifest.json` (still at the last
     stable version from `main`) should be kept as-is.
1. **Squash merge** the promotion PR to `main` with a conventional commit
   title that signals a breaking change, e.g.:

   ```
   feat!: release v5.0.0
   ```

   **Important:** GitHub's default squash message (`Merge pull request #123
   from v5-draft`) will be ignored by release-please. The merge author must
   write a proper `feat!:` or `fix!:` title.

1. release-please on `main`:
   - Reads the stable manifest (still at the last stable version, e.g., `4.0.0`)
   - Sees the new commit: `feat!: release v5.0.0`
   - Computes the next version: `4.0.0 + feat! → 5.0.0`
   - Creates a release PR for `v5.0.0`

1. A maintainer merges the release PR.
1. release-please creates the GitHub release with the `v5.0.0` tag.

**Why this works:** Because the promotion uses a **squash merge**, the RC
tags (which point to commits on the draft branch) are not in `main`'s commit
history. release-please's tag discovery scans only the last 250 commits on
the target branch, so it never sees the RC tags. The stable manifest
provides the explicit baseline, and the `feat!:` commit drives the natural
major version bump.

---

## Version Coordination & Gotchas

### Config Files on All Branches

Both `release-please-config.json` and `prerelease-please-config.json` exist
on all branches. The workflow selects the correct one at runtime. There is no
need for draft branches to create or modify their own config.

### Manifest Isolation

| Branch | Manifest Used | Typical Contents |
|--------|---------------|------------------|
| `main` | `.release-please-manifest.json` | `{ ".": "4.0.0" }` |
| `v5-draft` | `.prerelease-manifest.json` | `{ ".": "5.0.0-rc.3" }` |

Each manifest is branch-specific. The stable manifest on a draft branch is
stale (carried from when the branch was created) but preserved for promotion.

### go.mod on Draft Branches

The `go.mod` major version on the draft branch is different from `main`.
This is intentional — the validation workflow (`validate-release-pr.yml`)
verifies the match, and maintainers resolve it on merge.

### Changelog Overlap

Fixes cherry-picked or merged from `main` into the draft branch will appear
in both changelogs. This is acceptable — the same commit fixed the same issue
in both tracks.
