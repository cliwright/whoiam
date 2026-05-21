package internal

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const ExpectedEnvVar = "WHOIAM_EXPECTED_ENV"

type ConfigPath struct {
	Path string
	File string
}

type Config struct {
	Accounts map[string]string `yaml:"accounts"`
}

func NewConfigPath() (*ConfigPath, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &ConfigPath{Path: filepath.Join(homeDir, ".whoiam"), File: "whoiam.yaml"}, nil
}

// NewProjectConfigPath returns a ConfigPath rooted at .whoiam/ in the current directory.
// Used by whoiam init to create a project-local config.
func NewProjectConfigPath() (*ConfigPath, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &ConfigPath{Path: filepath.Join(dir, ".whoiam"), File: "whoiam.yaml"}, nil
}

// FindLocalConfigPath walks up the directory tree from CWD looking for a .whoiam/whoiam.yaml.
// Returns nil (no error) if no local config is found.
// Stops before reaching the home directory to avoid treating ~/.whoiam/whoiam.yaml as local.
func FindLocalConfigPath() (*ConfigPath, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	for {
		// Don't pick up the global config as a local one.
		if dir == homeDir {
			return nil, nil
		}
		candidate := filepath.Join(dir, ".whoiam", "whoiam.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return &ConfigPath{Path: filepath.Join(dir, ".whoiam"), File: "whoiam.yaml"}, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, nil
		}
		dir = parent
	}
}

// LoadEffectiveConfig loads the global config and merges any project-local config on top.
// Local account definitions take precedence over global ones on conflict.
func LoadEffectiveConfig() (*Config, error) {
	cfg, _, err := LoadEffectiveConfigWithSources()
	return cfg, err
}

// LoadEffectiveConfigWithSources is like LoadEffectiveConfig but also returns a map of
// account name -> source ("global" or "local") for each entry.
func LoadEffectiveConfigWithSources() (*Config, map[string]string, error) {
	globalPath, err := NewConfigPath()
	if err != nil {
		return nil, nil, err
	}

	cfg := &Config{Accounts: make(map[string]string)}
	sources := make(map[string]string)

	if globalPath.ConfigFileExists() {
		cfg, err = globalPath.LoadConfig()
		if err != nil {
			return nil, nil, err
		}
		for name := range cfg.Accounts {
			sources[name] = "global"
		}
	}

	localPath, err := FindLocalConfigPath()
	if err != nil {
		return nil, nil, err
	}

	if localPath == nil {
		return cfg, sources, nil
	}

	localCfg, err := localPath.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	for name, number := range localCfg.Accounts {
		cfg.Accounts[name] = number
		sources[name] = "local"
	}

	return cfg, sources, nil
}

// FindLocalDir walks up from CWD looking for a .whoiam/ directory.
// Returns "" (no error) if none is found.
// Stops before reaching the home directory to avoid treating ~/.whoiam/ as local.
func FindLocalDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	for {
		if dir == homeDir {
			return "", nil
		}
		candidate := filepath.Join(dir, ".whoiam")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

// ReadCurrentEnv reads the account name from .whoiam/expected-env.
// Checks local first, then falls back to global. Returns "" if neither is set.
func ReadCurrentEnv() (string, error) {
	name, _, err := ReadCurrentEnvWithSource()
	return name, err
}

// ReadCurrentEnvWithSource is like ReadCurrentEnv but also returns the source:
// "env" (WHOIAM_EXPECTED_ENV), "local" (.whoiam/expected-env), "global" (~/.whoiam/expected-env), or "" if not set.
func ReadCurrentEnvWithSource() (string, string, error) {
	if val := os.Getenv(ExpectedEnvVar); val != "" {
		return val, "env", nil
	}

	localDir, err := FindLocalDir()
	if err != nil {
		return "", "", err
	}
	if localDir != "" {
		data, err := os.ReadFile(filepath.Join(localDir, "expected-env"))
		if err == nil {
			name := strings.TrimSpace(string(data))
			if name != "" {
				return name, "local", nil
			}
		} else if !os.IsNotExist(err) {
			return "", "", err
		}
	}

	globalPath, err := NewConfigPath()
	if err != nil {
		return "", "", err
	}
	data, err := os.ReadFile(filepath.Join(globalPath.Path, "expected-env"))
	if os.IsNotExist(err) {
		return "", "", nil
	}
	if err != nil {
		return "", "", err
	}
	name := strings.TrimSpace(string(data))
	if name == "" {
		return "", "", nil
	}
	return name, "global", nil
}

// WriteCurrentEnv writes the account name to .whoiam/expected-env.
// Requires a .whoiam/ directory to already exist (created by whoiam init).
func WriteCurrentEnv(name string) error {
	localDir, err := FindLocalDir()
	if err != nil {
		return err
	}
	if localDir == "" {
		return fmt.Errorf("no .whoiam/ directory found — run 'whoiam init' first, or use --global to set globally")
	}
	return os.WriteFile(filepath.Join(localDir, "expected-env"), []byte(name+"\n"), 0644)
}

// ClearCurrentEnv removes .whoiam/expected-env if it exists.
func ClearCurrentEnv() error {
	localDir, err := FindLocalDir()
	if err != nil {
		return err
	}
	if localDir == "" {
		return nil
	}
	err = os.Remove(filepath.Join(localDir, "expected-env"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ReadGlobalCurrentEnv reads the account name from ~/.whoiam/expected-env.
func ReadGlobalCurrentEnv() (string, error) {
	globalPath, err := NewConfigPath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(globalPath.Path, "expected-env"))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// WriteGlobalCurrentEnv writes the account name to ~/.whoiam/expected-env.
func WriteGlobalCurrentEnv(name string) error {
	globalPath, err := NewConfigPath()
	if err != nil {
		return err
	}
	if !globalPath.Exists() {
		return fmt.Errorf("no ~/.whoiam/ directory found — run 'whoiam init --global' first")
	}
	return os.WriteFile(filepath.Join(globalPath.Path, "expected-env"), []byte(name+"\n"), 0644)
}

// ClearGlobalCurrentEnv removes ~/.whoiam/expected-env if it exists.
func ClearGlobalCurrentEnv() error {
	globalPath, err := NewConfigPath()
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(globalPath.Path, "expected-env"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func NewTemplateConfig() (*Config, error) {
	account := make(map[string]string)
	account["account"] = "123456789012"

	return &Config{
		Accounts: account,
	}, nil
}

func ValidateAccountNumber(number string) error {
	if len(number) != 12 {
		return fmt.Errorf("account number must be 12 digits")
	}

	if _, err := strconv.Atoi(number); err != nil {
		return fmt.Errorf("account number must only contain digits")
	}
	return nil
}

func (c *ConfigPath) Exists() bool {
	_, err := os.Stat(c.Path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *ConfigPath) ConfigFileExists() bool {
	_, err := os.Stat(c.FullPath())
	return err == nil
}

func (c *ConfigPath) FullPath() string {
	return filepath.Join(c.Path, c.File)
}

func (c *ConfigPath) Create() error {
	return os.MkdirAll(c.Path, 0700)
}

func (c *ConfigPath) LoadConfig() (*Config, error) {
	data, err := os.ReadFile(c.FullPath())
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *ConfigPath) SaveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(c.FullPath(), data, 0644)
}

func (c *Config) AccountExists(name string) bool {
	_, exists := c.Accounts[name]
	return exists
}

func (c *Config) AddAccount(name, number string) error {
	if err := ValidateAccountNumber(number); err != nil {
		return err
	}
	c.Accounts[name] = number
	return nil
}

func (c *Config) DeleteAccount(name string) {
	delete(c.Accounts, name)
}

func (c *Config) GetAccountByNumber(number string) string {
	for key, value := range c.Accounts {
		if value == number {
			return key
		}
	}
	return ""
}

func (c *Config) PrintConfigTable() {
	c.PrintConfigTableWithSource(nil)
}

func (c *Config) PrintConfigTableWithSource(sources map[string]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if sources != nil {
		table.SetHeader([]string{"Account Name", "Account Number", "Source"})
		for name, account := range c.Accounts {
			source := sources[name]
			table.Append([]string{name, account, source})
		}
	} else {
		table.SetHeader([]string{"Account Name", "Account Number"})
		for name, account := range c.Accounts {
			table.Append([]string{name, account})
		}
	}
	table.Render()
}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
