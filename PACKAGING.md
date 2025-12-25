# Taskr Package Manager Publishing Guide

This document provides step-by-step instructions for publishing Taskr to various package managers after a new release has been created.

## Prerequisites

- [x] Release created with: `taskr release:<version>`
- [x] GoReleaser CI completed: https://github.com/vikbert/taskr/actions/workflows/release.yml
- [x] Release tag created and pushed to GitHub
- [x] Release artifacts available on GitHub Releases

## Quick Commands Reference

```bash
# Check current release status
gh release list --repo vikbert/taskr

# Download release artifacts (if needed)
gh release download v3.46.3 --repo vikbert/taskr

# Update Homebrew formula example
brew bump-formula-pr --url https://github.com/vikbert/taskr/archive/v3.46.3.tar.gz taskr
```

## Official Package Managers

### 1. Homebrew (macOS/Linux)

**Status**: Requires manual setup
**Timeline**: After creating your own tap

**Steps:**
1. Create your own Homebrew tap: `vikbert/homebrew-taskr`
   ```bash
   # Fork the original tap as a starting point
   git clone https://github.com/go-task/homebrew-tap.git
   # Rename and modify for Taskr
   ```

2. Update `Formula/taskr.rb`:
   - Change URLs to point to `vikbert/taskr`
   - Update version and SHA256 hash
   - Update description to mention Taskr

3. Test locally:
   ```bash
   brew install --build-from-source taskr
   brew test taskr
   ```

4. Submit PR to your tap repository

### 2. npm (Cross-platform)

**Status**: Requires setup
**Timeline**: Immediate after release

**Steps:**
1. Publish to npm:
   ```bash
   cd packaging/npm
   npm publish
   ```

2. Package name options:
   - `@vikbert/taskr`
   - `@taskr/cli`

3. Ensure `package.json` has correct version and metadata

### 3. Snap (Linux)

**Status**: Repository exists
**Timeline**: Immediate after release

**Steps:**
1. Go to: https://github.com/vikbert/snap
2. Update `snap/snapcraft.yaml`:
   - Update version number
   - Update source URL to point to new release
   - Update changelog

3. Push changes to trigger automated build

### 4. WinGet (Windows)

**Status**: Requires community contribution
**Timeline**: After release

**Steps:**
1. Go to: https://github.com/microsoft/winget-pkgs
2. Find or create `Taskr.Task` package manifest
3. Update manifest with new version and installer URL
4. Submit pull request

## Cloudsmith (RPM/DEB)

**Status**: Requires setup
**Timeline**: After creating repository

**Steps:**
1. Create Taskr repository at: https://cloudsmith.io/~vikbert/repos/taskr
2. Upload RPM/DEB packages from GoReleaser artifacts
3. Update repository metadata
4. Update installation URLs in docs

## Community Package Managers

These require submitting pull requests to the respective package repositories:

### Chocolatey (Windows)
- Repository: https://github.com/chocolatey-community/chocolatey-coreteampackages
- Package name: `taskr`
- Submit PR with updated package definition

### Scoop (Windows)
- Repository: https://github.com/ScoopInstaller/Main
- Path: `bucket/task.json`
- Submit PR to update the manifest

### Arch Linux
- Repository: https://gitlab.archlinux.org/archlinux/packaging/packages/go-task
- Package name: `go-task` (consider requesting rename to `taskr`)
- Submit MR to update PKGBUILD

### Fedora
- Repository: https://src.fedoraproject.org/rpms/golang-github-task
- Package name: `golang-github-task` (consider requesting rename)
- Submit PR to update spec file

### Nix
- Repository: https://github.com/NixOS/nixpkgs
- Path: `pkgs/by-name/go/go-task/package.nix`
- Submit PR to update package definition

### FreeBSD
- Repository: https://cgit.freebsd.org/ports/tree/devel/task
- Path: `devel/task/`
- Submit PR to update Makefile and distinfo

## Automation Status

### ✅ Already Automated
- **GitHub Releases**: Via GoReleaser
- **Mise**: Configured for `vikbert/taskr`
- **pkgx**: Configured for `taskr-io.vercel.app`

### 🔄 Requires Manual Setup
- **GitHub Action**: Fork of `go-task/setup-task` needed
- **Homebrew Tap**: Custom tap needed
- **Cloudsmith**: Repository setup needed

## Post-Publishing Checklist

- [ ] Update installation documentation with new package URLs
- [ ] Test installations on all supported platforms
- [ ] Update README badges and links
- [ ] Announce release on social media/communities
- [ ] Monitor for issues and provide support

## Useful Links

- [Taskr Releases](https://github.com/vikbert/taskr/releases)
- [GoReleaser Artifacts](https://github.com/vikbert/taskr/actions/workflows/release.yml)
- [Original Task Packaging](https://github.com/go-task/task/blob/main/.goreleaser.yml)
- [Package Manager Guidelines](https://github.com/go-task/task/blob/main/docs/PACKAGING.md)

## Getting Help

If you need help with any of these packaging steps:

1. Check the original Task project's packaging documentation
2. Look at existing PRs in the target repositories
3. Ask in the Taskr community or GitHub Discussions
4. Reach out to maintainers of the target package repositories
