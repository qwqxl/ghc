# ghc - GitHub 配置管理工具

一个简化的 CLI 工具，用于管理 GitHub 仓库配置与项目版本，无需复杂的 Git 命令。

## 功能特性

- 🚀 **极简 CLI**：一条命令即可完成版本管理和仓库操作
- 📁 **配置管理**：自动生成和管理配置文件
- 🔒 **仓库绑定**：通过锁定文件避免操作失误
- 🏷️ **版本管理**：自动管理 Git 标签和版本号
- 📈 **历史追溯**：方便快速回滚到指定版本

## 安装

```bash
# 编译项目
go build -o ghc .

# 将 ghc 添加到系统 PATH 中使用
```

## 快速开始

### 1. 初始化项目

```bash
ghc init
```

这将创建：
- `ghc.config.yaml` - 项目配置文件
- `.repo.lock` - 仓库锁定文件

### 2. 绑定仓库

```bash
ghc bind https://github.com/username/project.git
```

### 3. 查看状态

```bash
ghc status
```

输出示例：
```
repo: https://github.com/username/project.git
branch: main
version: 0.0.1
tag_prefix: v
auto_push: true
build_command: go build ./...
```

### 4. 版本管理

```bash
# 创建新版本标签
ghc tag 1.0.0

# 查看所有标签
ghc tag list

# 切换到指定版本
ghc tag checkout 1.0.0
```

## 配置文件

### ghc.config.yaml

```yaml
repo: "https://github.com/username/project.git"  # 仓库地址
branch: main                                     # 默认分支
auto_push: true                                  # 自动推送
build_command: "go build ./..."                  # 构建命令
version: 0.0.1                                   # 当前版本
tag_prefix: v                                    # 标签前缀
```

### .repo.lock

```yaml
repo: https://github.com/username/project.git
branch: main
```

## 命令参考

| 命令 | 描述 |
|------|------|
| `ghc init` | 初始化项目配置 |
| `ghc bind <repo-url>` | 绑定仓库地址 |
| `ghc status` | 查看当前状态 |
| `ghc tag <version>` | 创建新标签 |
| `ghc tag list` | 查看所有标签 |
| `ghc tag checkout <version>` | 切换到指定版本 |
| `ghc help` | 显示帮助信息 |

## 开发

```bash
# 运行项目
go run .

# 编译
go build -o ghc .

# 测试
go test ./...
```

## 许可证

MIT License