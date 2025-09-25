# ghc - GitHub 配置管理工具（简化 CLI + 自动版本管理）

## ghc工具管理配置

* 在ghc工具还没有打包前 ghc项目自定义配置

```yaml
config_filename: ghc.config.yaml
repo_lock_filename: .repo.lock
default_repo: https://github.com/username/project.git
default_version: 0.0.1
default_branch: main
default_tag_prefix: v
cli:
  init: init
  bind: bind
  tag: tag
  st: status
  ...
```

## 1. 核心目标

* 管理 GitHub 仓库配置与项目版本，无需复杂 CLI 命令。
* 自动生成锁文件 `.repo.lock`，保持仓库绑定状态。
* 自动管理 Git 标签和版本号。
* CLI 命令简化到一个主命令 `ghc`，通过参数完成各种操作。

---

## 2. 配置文件与锁文件

### 2.1 `ghc.config.yaml` 示例

```yaml
repo: "https://github.com/username/project.git"  # 仓库地址
branch: main
auto_push: true
build_command: "go build ./..."
version: 0.0.1      # 当前版本
tag_prefix: v       # 标签前缀
```

### 2.2 `.repo.lock` 文件

* 记录仓库绑定状态和默认分支信息。
* 当 `.repo.lock` 存在时，表示仓库已经绑定；第一次初始化时自动生成。

```text
repo: https://github.com/username/project.git
branch: main
```

---

## 3. 仓库绑定逻辑

* 当 `ghc.config.yaml` 中 `repo` 字段为空：

```text
仓库未绑定，请使用 ghc bind <repo-url> 绑定仓库
```

* 第一次初始化（`ghc init`）时：

  1. 创建 `ghc.config.yaml`
  2. 创建 `.repo.lock`
  3. 生成默认版本 `0.0.1`

* `ghc bind <repo-url>`：绑定仓库地址并生成 `.repo.lock`

---

## 4. CLI 简化操作

* 所有操作都通过 `ghc` 主命令完成：

### 4.1 查看状态

```bash
ghc status
```

输出示例：

```
repo: https://github.com/username/project.git
branch: main
version: 1.0.0
tag_prefix: v
```

### 4.2 创建新标签（自动更新版本）

```bash
ghc tag 1.1.0
```

* 自动打标签 `v1.1.0`（根据 `tag_prefix` 配置）
* 更新 `ghc.config.yaml` 的 `version` 字段
* 自动 push 到远程仓库（如果 `auto_push: true`）

### 4.3 查看已有标签

```bash
ghc tag list
```

输出示例：

```
v1.0.0
v1.1.0
v1.1.1
```

### 4.4 切换到指定版本

```bash
ghc tag checkout 1.0.0
```

* 自动切换分支到对应标签 commit

### 4.5 绑定仓库

```bash
ghc bind https://github.com/username/project.git
```

* 更新 `ghc.config.yaml` 的 `repo` 字段
* 生成 `.repo.lock`

---

## 5. 操作流程（整合版本管理与仓库绑定）

1. **初始化项目**

```bash
ghc init
```

* 创建 `ghc.config.yaml`
* 生成 `.repo.lock`
* 默认版本 `0.0.1`

2. **绑定仓库（如果未绑定）**

```bash
ghc bind <repo-url>
```

3. **查看状态**

```bash
ghc status
```

4. **发布新版本**

```bash
ghc tag 1.1.0
```

5. **查看历史版本**

```bash
ghc tag list
```

6. **回退到某个版本**

```bash
ghc tag checkout 1.0.0
```

---

## 6. 技术实现建议

* **语言**：Go
* **库**：`go-git` 管理 Git 操作（标签、分支、push）
* **配置文件**：`ghc.config.yaml`
* **锁文件**：`.repo.lock` 用于标记仓库绑定状态
* **CLI 设计**：只需一个主命令 `ghc`，通过参数完成所有操作
* **自动化**：支持自动更新版本、自动 push、版本回退

---

## 7. 产品价值

* **极简 CLI**：一条命令即可完成版本管理和仓库操作
* **仓库绑定可追踪**：通过 `.repo.lock` 文件避免操作失误
* **自动版本管理**：创建标签、更新版本、push 自动完成
* **历史版本可追溯**：方便快速回滚到指定版本