package main

import (
  "fmt"
  "os"
)

func main() {
  if len(os.Args) < 2 {
    showHelp()
    return
  }

  command := os.Args[1]
  args := os.Args[2:]

  switch command {
  case "init":
    handleInit()
  case "bind":
    handleBind(args)
  case "status":
    handleStatus()
  case "tag":
    handleTag(args)
  case "publish", "release":
    handlePublish(args)
  case "help", "-h", "--help":
    showHelp()
  default:
    fmt.Printf("未知命令: %s\n", command)
    showHelp()
  }
}

func showHelp() {
  fmt.Println("ghc - GitHub 配置管理工具")
  fmt.Println("")
  fmt.Println("使用方法:")
  fmt.Println("  ghc init                    初始化项目配置")
  fmt.Println("  ghc bind <repo-url>         绑定仓库地址")
  fmt.Println("  ghc status                  查看当前状态")
  fmt.Println("  ghc tag <version>           创建新标签")
  fmt.Println("  ghc tag list                查看所有标签")
  fmt.Println("  ghc tag checkout <version>  切换到指定版本")
  fmt.Println("  ghc publish [version]       发布项目到 GitHub")
  fmt.Println("  ghc release [version]       发布项目到 GitHub (同 publish)")
  fmt.Println("  ghc help                    显示帮助信息")
}