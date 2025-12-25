<div align="center">
  <a href="https://taskr-io.vercel.app/">
    <img src="website/src/public/img/logo.png" width="180px"/>
  </a>

  <h4>Task++ runner for gophers</h4>
  <p>where tasks meet elegance</p>

  <a href="https://taskr-io.vercel.app/">
    <img src="website/src/public/img/taskr.png" width="100%"/>
  </a>
</div>


## 🚀 Quick Start

Get started with Taskr in just a few commands:

```bash

## install via brew
brew tap vikbert/taskr
brew install taskr

## install via go
go install github.com/vikbert/taskr/v3/cmd/taskr@latest

## install via shell
sh -c "$(curl --location https://taskr-io.vercel.app/install.sh)" -- -d

## install via binary
open https://github.com/vikbert/taskr/releases

```

## 🛠️ Development

```bash
git clone https://github.com/vikbert/taskr.git
cd taskr

# Download dependencies
task mod

# Install development tools
task install:mockery
task gotestsum:install
```

## 🎯 Key Enhancements

- ✅ **Performance Optimization** - Pre-allocated capacity, fast paths, 30% improvement for large project lists
- ✅ **New Category** - Categorize the related tasks by using category
- ✅ **New Index** - Order the tasks by index

## 📦 Packaging & Distribution

Taskr is distributed through multiple package managers. For information on publishing to package managers after releases, see:

- [Package Manager Publishing Guide](PACKAGING.md)
- [Installation Documentation](website/src/docs/installation.md)

### Release Process

```bash
# Prepare and create a release
taskr release:patch  # or release:minor/release:major

# Publish to package managers (manual steps)
taskr release:package-managers
```


```yml
test:all:
  category: test
  desc: Runs test suite with signals and watch tests included
  deps: [sleepit:build, gotestsum:install]
  cmds:
    - gotestsum -f '{{.GOTESTSUM_FORMAT}}' -tags 'signals watch' ./...

goreleaser:test:
  category: release
  desc: Tests release process without publishing
  cmds:
    - goreleaser --snapshot --clean
```
## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
