package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

// handlePublish 处理发布命令
func handlePublish(args []string) {
	// 检查是否请求帮助
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help") {
		fmt.Println("ghc publish/release - 发布项目到 GitHub")
		fmt.Println("")
		fmt.Println("使用方法:")
		fmt.Println("  ghc publish [version]    发布项目到 GitHub")
		fmt.Println("  ghc release [version]    发布项目到 GitHub (同 publish)")
		fmt.Println("")
		fmt.Println("参数:")
		fmt.Println("  version                  发布版本号 (可选，默认使用配置文件中的版本)")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Println("  ghc publish v1.0.0       发布版本 v1.0.0")
		fmt.Println("  ghc release              使用配置文件中的版本发布")
		return
	}

	// 获取版本号参数
	var version string
	if len(args) > 0 {
		version = args[0]
	} else {
		// 如果没有提供版本号，尝试从配置文件获取
		config, err := LoadConfig()
		if err != nil {
			fmt.Printf("加载配置失败: %v\n", err)
			return
		}
		version = config.Version
		if version == "" {
			version = "v1.0.0" // 默认版本
		}
	}

	fmt.Printf("开始发布项目，版本: %s\n", version)

	// 1. 编译项目
	fmt.Println("步骤 1/6: 编译项目...")
	if err := buildProject(); err != nil {
		fmt.Printf("编译失败: %v\n", err)
		return
	}
	fmt.Println("✓ 编译成功")

	// 2. 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}

	// 3. 初始化 Git 仓库（如果需要）
	fmt.Println("步骤 2/6: 检查 Git 仓库...")
	if !IsGitRepository(cwd) {
		fmt.Println("初始化 Git 仓库...")
		if err := InitRepository(cwd); err != nil {
			fmt.Printf("初始化 Git 仓库失败: %v\n", err)
			return
		}
	}
	fmt.Println("✓ Git 仓库就绪")

	// 4. 添加远程仓库
	fmt.Println("步骤 3/6: 配置远程仓库...")
	if err := setupRemoteRepository(); err != nil {
		fmt.Printf("配置远程仓库失败: %v\n", err)
		return
	}
	fmt.Println("✓ 远程仓库配置完成")

	// 5. 提交所有文件
	fmt.Println("步骤 4/6: 提交文件...")
	if err := commitAllFiles(version); err != nil {
		fmt.Printf("提交文件失败: %v\n", err)
		return
	}
	fmt.Println("✓ 文件提交完成")

	// 6. 推送到 GitHub
	fmt.Println("步骤 5/6: 推送到 GitHub...")
	if err := pushToGitHub(); err != nil {
		fmt.Printf("推送失败: %v\n", err)
		return
	}
	fmt.Println("✓ 推送完成")

	// 7. 创建发布标签
	fmt.Println("步骤 6/6: 创建发布标签...")
	if err := createReleaseTag(version); err != nil {
		fmt.Printf("创建标签失败: %v\n", err)
		return
	}
	fmt.Println("✓ 发布标签创建完成")

	fmt.Printf("\n🎉 项目发布成功！版本: %s\n", version)
}

// buildProject 编译项目
func buildProject() error {
	config, err := LoadConfig()
	if err != nil {
		// 如果没有配置文件，使用默认构建命令
		return runCommand("go build ./...")
	}

	// 执行预编译钩子
	if err := executePreBuildHooks(config); err != nil {
		return fmt.Errorf("预编译失败: %v", err)
	}

	// 执行主构建命令
	if config.BuildCommand == "" {
		return runCommand("go build ./...")
	}

	return runCommand(config.BuildCommand)
}

// setupRemoteRepository 设置远程仓库
func setupRemoteRepository() error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	if config.Repo == "" {
		return fmt.Errorf("未配置远程仓库地址，请先使用 'ghc bind <repo-url>' 绑定仓库")
	}

	// 检查是否已经添加了远程仓库
	cwd, _ := os.Getwd()
	gitOps, err := NewGitOperations(cwd)
	if err == nil {
		// 如果能获取到远程 URL，说明已经配置了
		if _, err := gitOps.GetRemoteURL(); err == nil {
			return nil // 远程仓库已存在
		}
	}

	// 添加远程仓库
	return runCommand(fmt.Sprintf("git remote add origin %s", config.Repo))
}

// commitAllFiles 提交所有文件
func commitAllFiles(version string) error {
	// 添加所有文件
	if err := runCommand("git add ."); err != nil {
		return fmt.Errorf("添加文件失败: %v", err)
	}

	// 提交文件 - 使用 exec.Command 直接处理参数
	commitMessage := fmt.Sprintf("Release version %s", version)
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// pushToGitHub 推送到 GitHub
func pushToGitHub() error {
	// 获取当前分支名
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	gitOps, err := NewGitOperations(cwd)
	if err != nil {
		return fmt.Errorf("初始化 Git 操作失败: %v", err)
	}

	branch, err := gitOps.GetCurrentBranch()
	if err != nil {
		// 如果获取失败，使用配置文件中的分支或默认分支
		config, configErr := LoadConfig()
		if configErr == nil && config.Branch != "" {
			branch = config.Branch
		} else {
			branch = "main"
		}
	}

	return runCommand(fmt.Sprintf("git push -u origin %s", branch))
}

// createReleaseTag 创建发布标签
func createReleaseTag(version string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	gitOps, err := NewGitOperations(cwd)
	if err != nil {
		return fmt.Errorf("初始化 Git 操作失败: %v", err)
	}

	// 创建标签
	tagMessage := fmt.Sprintf("Release version %s", version)
	if err := gitOps.CreateTag(version, tagMessage); err != nil {
		return fmt.Errorf("创建标签失败: %v", err)
	}

	// 推送标签
	if err := gitOps.PushTag(version); err != nil {
		return fmt.Errorf("推送标签失败: %v", err)
	}

	// 更新配置文件中的版本号
	config, err := LoadConfig()
	if err == nil {
		config.Version = version
		SaveConfig(config)
	}

	return nil
}

// executePreBuildHooks 执行预编译钩子
func executePreBuildHooks(config *Config) error {
	if !config.PreBuild.Enabled {
		return nil // 预编译未启用，直接返回
	}

	fmt.Println("🔧 执行预编译钩子...")

	// 执行预编译脚本（如果有）
	if config.PreBuild.Script != "" {
		fmt.Printf("执行预编译脚本: %s\n", config.PreBuild.Script)
		if err := runCommandWithTimeout(config.PreBuild.Script, config.PreBuild.Timeout); err != nil {
			if config.PreBuild.FailOnError {
				return fmt.Errorf("预编译脚本执行失败: %v", err)
			}
			fmt.Printf("⚠️ 预编译脚本执行失败（已忽略）: %v\n", err)
		}
	}

	// 执行预编译命令列表
	for i, command := range config.PreBuild.Commands {
		if command == "" {
			continue
		}
		fmt.Printf("执行预编译命令 [%d/%d]: %s\n", i+1, len(config.PreBuild.Commands), command)
		if err := runCommandWithTimeout(command, config.PreBuild.Timeout); err != nil {
			if config.PreBuild.FailOnError {
				return fmt.Errorf("预编译命令执行失败: %v", err)
			}
			fmt.Printf("⚠️ 预编译命令执行失败（已忽略）: %v\n", err)
		}
	}

	fmt.Println("✓ 预编译钩子执行完成")
	return nil
}

// runCommandWithTimeout 执行带超时的系统命令
func runCommandWithTimeout(command string, timeoutSeconds int) error {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 300 // 默认5分钟超时
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// 分割命令和参数
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("空命令")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("命令执行超时（%d秒）", timeoutSeconds)
	}
	return err
}

// runCommand 执行系统命令
func runCommand(command string) error {
	fmt.Printf("执行命令: %s\n", command)

	// 分割命令和参数
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("空命令")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
