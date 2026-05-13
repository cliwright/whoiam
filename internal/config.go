package internal

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type ConfigPath struct {
	Path string
	File string
}

type Config struct {
	Accounts map[string]string `yaml:"accounts"`
}

func NewConfigPath() (*ConfigPath, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	path := usr.HomeDir + "/.whoiam"
	return &ConfigPath{Path: path, File: "whoiam.yaml"}, nil
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
func FindLocalConfigPath() (*ConfigPath, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
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
	globalPath, err := NewConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := &Config{Accounts: make(map[string]string)}
	if globalPath.ConfigFileExists() {
		cfg, err = globalPath.LoadConfig()
		if err != nil {
			return nil, err
		}
	}

	localPath, err := FindLocalConfigPath()
	if err != nil {
		return nil, err
	}

	if localPath == nil {
		return cfg, nil
	}

	localCfg, err := localPath.LoadConfig()
	if err != nil {
		return nil, err
	}

	for name, number := range localCfg.Accounts {
		cfg.Accounts[name] = number
	}

	return cfg, nil
}

// FindLocalDir walks up from CWD looking for a .whoiam/ directory.
// Returns "" (no error) if none is found.
func FindLocalDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
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

// ReadCurrentEnv reads the account name from .whoiam/current-env.
// Returns "" (no error) if no current-env file exists.
func ReadCurrentEnv() (string, error) {
	localDir, err := FindLocalDir()
	if err != nil {
		return "", err
	}
	if localDir == "" {
		return "", nil
	}
	data, err := os.ReadFile(filepath.Join(localDir, "current-env"))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// WriteCurrentEnv writes the account name to .whoiam/current-env.
// Requires a .whoiam/ directory to already exist (created by whoiam init).
func WriteCurrentEnv(name string) error {
	localDir, err := FindLocalDir()
	if err != nil {
		return err
	}
	if localDir == "" {
		return fmt.Errorf("no .whoiam/ directory found — run 'whoiam init' first")
	}
	return os.WriteFile(filepath.Join(localDir, "current-env"), []byte(name+"\n"), 0644)
}

// ClearCurrentEnv removes .whoiam/current-env if it exists.
func ClearCurrentEnv() error {
	localDir, err := FindLocalDir()
	if err != nil {
		return err
	}
	if localDir == "" {
		return nil
	}
	err = os.Remove(filepath.Join(localDir, "current-env"))
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
	return c.Path + "/" + c.File
}

func (c *ConfigPath) Create() error {
	err := os.MkdirAll(c.Path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
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

	err = os.WriteFile(c.FullPath(), data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Account Name", "Account Number"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	for name, account := range c.Accounts {
		table.Append([]string{name, account})
	}
	table.Render()
}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
