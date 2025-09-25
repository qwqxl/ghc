package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// PreBuildConfig 预编译配置结构
type PreBuildConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Commands    []string `yaml:"commands"`
	Script      string   `yaml:"script"`
	Timeout     int      `yaml:"timeout"`
	FailOnError bool     `yaml:"fail_on_error"`
}

// Config 项目配置结构
type Config struct {
	Repo         string         `yaml:"repo"`
	Branch       string         `yaml:"branch"`
	AutoPush     bool           `yaml:"auto_push"`
	BuildCommand string         `yaml:"build_command"`
	Version      string         `yaml:"version"`
	TagPrefix    string         `yaml:"tag_prefix"`
	PreBuild     PreBuildConfig `yaml:"pre_build"`
}

// RepoLock 仓库锁定文件结构
type RepoLock struct {
	Repo           string `yaml:"repo"`
	Branch         string `yaml:"branch"`
	CurrentVersion string `yaml:"current_version"`
	LastUpdated    string `yaml:"last_updated"`
}

const (
	ConfigFile   = "ghc.config.yaml"
	RepoLockFile = ".repo.lock"
)

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	if !fileExists(ConfigFile) {
		return nil, fmt.Errorf("配置文件 %s 不存在，请先运行 ghc init", ConfigFile)
	}

	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// SaveConfig 保存配置文件
func SaveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	err = ioutil.WriteFile(ConfigFile, data, 0644)
	if err != nil {
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	return nil
}

// LoadRepoLock 加载仓库锁定文件
func LoadRepoLock() (*RepoLock, error) {
	if !fileExists(RepoLockFile) {
		return nil, fmt.Errorf("仓库锁定文件 %s 不存在", RepoLockFile)
	}

	data, err := ioutil.ReadFile(RepoLockFile)
	if err != nil {
		return nil, fmt.Errorf("读取仓库锁定文件失败: %v", err)
	}

	var lock RepoLock
	err = yaml.Unmarshal(data, &lock)
	if err != nil {
		return nil, fmt.Errorf("解析仓库锁定文件失败: %v", err)
	}

	return &lock, nil
}

// SaveRepoLock 保存仓库锁定文件
func SaveRepoLock(lock *RepoLock) error {
	data, err := yaml.Marshal(lock)
	if err != nil {
		return fmt.Errorf("序列化仓库锁定信息失败: %v", err)
	}

	err = ioutil.WriteFile(RepoLockFile, data, 0644)
	if err != nil {
		return fmt.Errorf("保存仓库锁定文件失败: %v", err)
	}

	return nil
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// loadRepoLock 加载仓库锁定文件（小写别名）
func loadRepoLock() (*RepoLock, error) {
	return LoadRepoLock()
}

// saveRepoLock 保存仓库锁定文件（小写别名）
func saveRepoLock(lock *RepoLock) error {
	return SaveRepoLock(lock)
}
