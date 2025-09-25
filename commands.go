package main

import (
  "fmt"
  "os"
  "strings"
  "time"
)

// handleInit 处理初始化命令
func handleInit() {
  if fileExists(ConfigFile) {
    fmt.Printf("项目已经初始化，配置文件 %s 已存在\n", ConfigFile)
    return
  }

  // 创建默认配置
  config := &Config{
    Repo:         "",
    Branch:       "main",
    AutoPush:     true,
    BuildCommand: "go build ./...",
    Version:      "0.0.1",
    TagPrefix:    "v",
  }

  err := SaveConfig(config)
  if err != nil {
    fmt.Printf("创建配置文件失败: %v\n", err)
    return
  }

  // 创建仓库锁定文件
  lock := &RepoLock{
    Repo:   "",
    Branch: "main",
  }

  err = SaveRepoLock(lock)
  if err != nil {
    fmt.Printf("创建仓库锁定文件失败: %v\n", err)
    return
  }

  fmt.Println("项目初始化成功！")
  fmt.Printf("已创建配置文件: %s\n", ConfigFile)
  fmt.Printf("已创建锁定文件: %s\n", RepoLockFile)
  fmt.Println("请使用 'ghc bind <repo-url>' 绑定仓库")
}

// handleBind 处理仓库绑定命令
func handleBind(args []string) {
  if len(args) == 0 {
    fmt.Println("请提供仓库地址")
    fmt.Println("使用方法: ghc bind <repo-url>")
    return
  }

  repoUrl := args[0]
  if !strings.HasPrefix(repoUrl, "https://github.com/") && !strings.HasPrefix(repoUrl, "git@github.com:") {
    fmt.Println("请提供有效的 GitHub 仓库地址")
    return
  }

  // 加载配置文件
  config, err := LoadConfig()
  if err != nil {
    fmt.Printf("加载配置失败: %v\n", err)
    return
  }

  // 更新配置
  config.Repo = repoUrl
  err = SaveConfig(config)
  if err != nil {
    fmt.Printf("保存配置失败: %v\n", err)
    return
  }

  // 更新锁定文件
  lock := &RepoLock{
    Repo:   repoUrl,
    Branch: config.Branch,
  }
  err = SaveRepoLock(lock)
  if err != nil {
    fmt.Printf("保存仓库锁定文件失败: %v\n", err)
    return
  }

  fmt.Printf("仓库绑定成功: %s\n", repoUrl)
}

// handleStatus 处理状态查看命令
func handleStatus() {
  config, err := LoadConfig()
  if err != nil {
    fmt.Printf("加载配置失败: %v\n", err)
    return
  }

  if config.Repo == "" {
    fmt.Println("仓库未绑定，请使用 ghc bind <repo-url> 绑定仓库")
    return
  }

  fmt.Printf("repo: %s\n", config.Repo)
  fmt.Printf("branch: %s\n", config.Branch)
  fmt.Printf("version: %s\n", config.Version)
  fmt.Printf("tag_prefix: %s\n", config.TagPrefix)
  fmt.Printf("auto_push: %t\n", config.AutoPush)
  fmt.Printf("build_command: %s\n", config.BuildCommand)
}

// handleTag 处理标签相关命令
func handleTag(args []string) {
  if len(args) == 0 {
    fmt.Println("请提供标签操作参数")
    fmt.Println("使用方法:")
    fmt.Println("  ghc tag <version>           创建新标签")
    fmt.Println("  ghc tag list                查看所有标签")
    fmt.Println("  ghc tag checkout <version>  切换到指定版本")
    return
  }

  subCommand := args[0]
  switch subCommand {
  case "list":
    handleTagList()
  case "checkout":
    if len(args) < 2 {
      fmt.Println("请提供要切换的版本号")
      return
    }
    handleTagCheckout(args[1])
  default:
    // 默认为创建标签
    handleTagCreate(subCommand)
  }
}

// handleTagCreate 创建新标签
func handleTagCreate(version string) {
	// 验证版本号格式
	if version == "" {
		fmt.Println("Error: Version cannot be empty")
		return
	}

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// 检查是否为 Git 仓库
	if !IsGitRepository(cwd) {
		fmt.Println("Error: Not a git repository")
		return
	}

	// 创建 Git 操作实例
	gitOps, err := NewGitOperations(cwd)
	if err != nil {
		fmt.Printf("Error initializing git operations: %v\n", err)
		return
	}

	// 验证仓库状态
	if err := gitOps.ValidateRepository(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 创建标签
	tagMessage := fmt.Sprintf("Release version %s", version)
	if err := gitOps.CreateTag(version, tagMessage); err != nil {
		fmt.Printf("Error creating tag: %v\n", err)
		return
	}

	// 推送标签到远程仓库
	if err := gitOps.PushTag(version); err != nil {
		fmt.Printf("Error pushing tag: %v\n", err)
		return
	}

	// 更新 .repo.lock 文件
	repoLock, err := loadRepoLock()
	if err != nil {
		fmt.Printf("Warning: Could not load repo lock: %v\n", err)
	} else {
		repoLock.CurrentVersion = version
		repoLock.LastUpdated = time.Now().Format(time.RFC3339)
		if err := saveRepoLock(repoLock); err != nil {
			fmt.Printf("Warning: Could not save repo lock: %v\n", err)
		}
	}

	fmt.Printf("Tag '%s' created and pushed successfully\n", version)
}

// handleTagList 列出所有标签
func handleTagList() {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// 检查是否为 Git 仓库
	if !IsGitRepository(cwd) {
		fmt.Println("Error: Not a git repository")
		return
	}

	// 创建 Git 操作实例
	gitOps, err := NewGitOperations(cwd)
	if err != nil {
		fmt.Printf("Error initializing git operations: %v\n", err)
		return
	}

	// 获取标签列表
	tags, err := gitOps.ListTags()
	if err != nil {
		fmt.Printf("Error listing tags: %v\n", err)
		return
	}

	if len(tags) == 0 {
		fmt.Println("No tags found")
		return
	}

	fmt.Println("Available tags:")
	for _, tag := range tags {
		fmt.Printf("  %s\n", tag)
	}

	// 显示当前标签
	latestTag, err := gitOps.GetLatestTag()
	if err == nil {
		fmt.Printf("\nLatest tag: %s\n", latestTag)
	}
}

// handleTagCheckout 切换到指定版本
func handleTagCheckout(version string) {
	// 验证版本号
	if version == "" {
		fmt.Println("Error: Version cannot be empty")
		return
	}

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// 检查是否为 Git 仓库
	if !IsGitRepository(cwd) {
		fmt.Println("Error: Not a git repository")
		return
	}

	// 创建 Git 操作实例
	gitOps, err := NewGitOperations(cwd)
	if err != nil {
		fmt.Printf("Error initializing git operations: %v\n", err)
		return
	}

	// 验证仓库状态
	if err := gitOps.ValidateRepository(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 切换到指定标签
	if err := gitOps.CheckoutTag(version); err != nil {
		fmt.Printf("Error checking out tag: %v\n", err)
		return
	}

	// 更新 .repo.lock 文件
	repoLock, err := loadRepoLock()
	if err != nil {
		fmt.Printf("Warning: Could not load repo lock: %v\n", err)
	} else {
		repoLock.CurrentVersion = version
		repoLock.LastUpdated = time.Now().Format(time.RFC3339)
		if err := saveRepoLock(repoLock); err != nil {
			fmt.Printf("Warning: Could not save repo lock: %v\n", err)
		}
	}

	fmt.Printf("Successfully checked out tag '%s'\n", version)
}