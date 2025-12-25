---
title: Installation
description: Installation methods for Taskr
outline: deep
---

# Installation

Taskr offers several installation methods. Check out the available methods
below.

## Official Package Managers

These installation methods are maintained by the Taskr team and are always
up-to-date.

::: info Current Status
- ✅ **Homebrew**: Available via `vikbert/taskr` tap
- 🔄 **Others**: Additional package managers (Snap, npm, WinGet, etc.) will be available in future releases
:::

### [Homebrew](https://brew.sh) ![macOS](https://img.shields.io/badge/MacOS-000000?logo=apple&logoColor=F0F0F0) ![Linux](https://img.shields.io/badge/Linux-FCC624?logo=linux&logoColor=black) {#homebrew}

Taskr is available via our official Homebrew tap:

```shell
brew tap vikbert/taskr
brew install taskr
```

[[source](https://github.com/vikbert/homebrew-taskr/blob/main/Formula/taskr.rb)]


## Go Install
```shell 
go install github.com/vikbert/taskr/v3/cmd/taskr@latest
```


## Binary Install

You can download the binary from the
[releases page on GitHub](https://github.com/vikbert/taskr/releases) and add to
your `$PATH`.

The `taskr_checksums.txt` file contains the SHA-256 checksum for each file.

## Shell Install

We also have an
[install script](https://github.com/vikbert/taskr/blob/main/install-task.sh)
which is very useful in scenarios like CI. Many thanks to
[GoDownloader](https://github.com/goreleaser/godownloader) for enabling the easy
generation of this script.

By default, it installs on the `./bin` directory relative to the working
directory:

```shell
sh -c "$(curl --location https://taskr-io.vercel.app/install.sh)" -- -d
```

::: tip Note
The install script requires release binaries to be available. If you just released a new version, please wait 5-15 minutes for our CI to build and upload the binaries.
:::

It is possible to override the installation directory with the `-b` parameter.
On Linux, common choices are `~/.local/bin` and `~/bin` to install for the
current user or `/usr/local/bin` to install for all users:

```shell
sh -c "$(curl --location https://taskr-io.vercel.app/install.sh)" -- -d -b ~/.local/bin
```

::: warning

On macOS and Windows, `~/.local/bin` and `~/bin` are not added to `$PATH` by
default.

:::

By default, it installs the latest version available. You can also specify a tag
(available in [releases](https://github.com/vikbert/taskr/releases)) to install
a specific version:

```shell
sh -c "$(curl --location https://taskr-io.vercel.app/install.sh)" -- -d v3.36.0
```

Parameters are order specific, to set both installation directory and version:

```shell
sh -c "$(curl --location https://taskr-io.vercel.app/install.sh)" -- -d -b ~/.local/bin v3.42.1
```

