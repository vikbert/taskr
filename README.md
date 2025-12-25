<div align="center">
  <a href="https://taskr-io.vercel.app/">
    <img src="website/src/public/img/intro.png" width="100%"/>
  </a>

</div>


## 🚀 Quick Start

Get started with Taskr in just a few commands:

#### install via brew
```bash
brew tap vikbert/taskr
brew install taskr
```

#### install via go
```
go install github.com/vikbert/taskr/v3/cmd/taskr@latest
```

#### install via shell
```
sh -c "$(curl --location https://taskr-io.vercel.app/install.sh)" -- -d
```

#### install via binary
```
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

  <a href="https://taskr-io.vercel.app/">
    <img src="website/src/public/img/taskr.png" width="100%"/>
  </a>

- ✅ **Performance Optimization** - Pre-allocated capacity, fast paths, 30% improvement for large project lists
- ✅ **New Category** - Categorize the related tasks by using category
- ✅ **New Index** - Order the tasks by index
- ✅ **New banner printing** - Print project name as project banner

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
