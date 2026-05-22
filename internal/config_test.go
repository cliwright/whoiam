package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// chdirTemp changes the working directory to dir for the duration of the test,
// restoring the original directory when the test completes.
func chdirTemp(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

// newTempConfigPath creates a ConfigPath rooted in a temp dir (not ~/.whoiam).
func newTempConfigPath(t *testing.T) *ConfigPath {
	t.Helper()
	dir := t.TempDir()
	return &ConfigPath{Path: filepath.Join(dir, ".whoiam"), File: "whoiam.yaml"}
}

// --- ConfigPath ---

func TestConfigPathFullPath(t *testing.T) {
	cp := newTempConfigPath(t)
	expected := filepath.Join(cp.Path, cp.File)
	if got := cp.FullPath(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestConfigPathExists(t *testing.T) {
	cp := newTempConfigPath(t)
	if cp.Exists() {
		t.Error("expected false before creation")
	}
	cp.Create()
	if !cp.Exists() {
		t.Error("expected true after creation")
	}
}

func TestConfigPathCreate(t *testing.T) {
	cp := newTempConfigPath(t)
	if err := cp.Create(); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !cp.Exists() {
		t.Error("expected dir to exist after Create")
	}
	// Create is idempotent
	if err := cp.Create(); err != nil {
		t.Fatalf("Create (idempotent): %v", err)
	}
}

func TestConfigFileExists(t *testing.T) {
	cp := newTempConfigPath(t)
	if cp.ConfigFileExists() {
		t.Error("expected false before file is written")
	}
	cp.Create()
	cfg, _ := NewTemplateConfig()
	cp.SaveConfig(cfg)
	if !cp.ConfigFileExists() {
		t.Error("expected true after file is written")
	}
}

func TestConfigPathSaveAndLoadConfig(t *testing.T) {
	cp := newTempConfigPath(t)
	cp.Create()

	cfg, err := NewTemplateConfig()
	if err != nil {
		t.Fatal(err)
	}
	if err := cp.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	loaded, err := cp.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if len(loaded.Accounts) != 1 {
		t.Errorf("expected 1 account, got %d", len(loaded.Accounts))
	}
}

func TestConfigPathLoadConfig_FileNotFound(t *testing.T) {
	cp := newTempConfigPath(t)
	_, err := cp.LoadConfig()
	if err == nil {
		t.Error("expected error when file does not exist")
	}
}

// --- NewProjectConfigPath ---

func TestNewProjectConfigPath(t *testing.T) {
	dir := t.TempDir()
	chdirTemp(t, dir)

	cp, err := NewProjectConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	// Resolve symlinks (macOS /var/folders → /private/var/folders)
	wantPath, _ := filepath.EvalSymlinks(filepath.Join(dir, ".whoiam"))
	gotPath, _ := filepath.EvalSymlinks(cp.Path)
	if gotPath != wantPath {
		t.Errorf("expected path %q, got %q", wantPath, gotPath)
	}
	if cp.File != "whoiam.yaml" {
		t.Errorf("expected whoiam.yaml, got %q", cp.File)
	}
}

// --- ValidateAccountNumber ---

func TestValidateAccountNumber(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"123456789012", false},
		{"000000000000", false},
		{"12345678901", true},    // 11 digits — too short
		{"1234567890123", true},  // 13 digits — too long
		{"12345678901a", true},   // contains letter
		{"", true},               // empty
	}
	for _, tt := range tests {
		err := ValidateAccountNumber(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateAccountNumber(%q): got err=%v, wantErr=%v", tt.input, err, tt.wantErr)
		}
	}
}

// --- Config methods ---

func TestConfigAccountExists(t *testing.T) {
	cfg, _ := NewTemplateConfig()
	if !cfg.AccountExists("account") {
		t.Error("expected true for existing account")
	}
	if cfg.AccountExists("nonexistent") {
		t.Error("expected false for nonexistent account")
	}
}

func TestConfigAddAccount(t *testing.T) {
	cfg, _ := NewTemplateConfig()
	if err := cfg.AddAccount("newaccount", "123456789012"); err != nil {
		t.Fatalf("AddAccount: %v", err)
	}
	if !cfg.AccountExists("newaccount") {
		t.Error("expected account to exist after add")
	}
}

func TestConfigAddAccount_InvalidNumber(t *testing.T) {
	cfg, _ := NewTemplateConfig()
	if err := cfg.AddAccount("bad", "123"); err == nil {
		t.Error("expected error for invalid account number")
	}
}

func TestConfigDeleteAccount(t *testing.T) {
	cfg, _ := NewTemplateConfig()
	cfg.DeleteAccount("account")
	if cfg.AccountExists("account") {
		t.Error("expected false after delete")
	}
}

func TestConfigGetAccountByNumber(t *testing.T) {
	cfg, _ := NewTemplateConfig()
	if got := cfg.GetAccountByNumber("123456789012"); got != "account" {
		t.Errorf("expected account, got %q", got)
	}
	if got := cfg.GetAccountByNumber("000000000000"); got != "" {
		t.Errorf("expected empty string for unknown number, got %q", got)
	}
}

// --- FindLocalDir ---

func TestFindLocalDir_Found(t *testing.T) {
	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	sub := filepath.Join(root, "nested", "subdir")
	os.MkdirAll(whoiamDir, 0700)
	os.MkdirAll(sub, 0700)
	chdirTemp(t, sub)

	found, err := FindLocalDir()
	if err != nil {
		t.Fatal(err)
	}
	// Resolve symlinks (macOS /var/folders → /private/var/folders)
	wantDir, _ := filepath.EvalSymlinks(whoiamDir)
	gotDir, _ := filepath.EvalSymlinks(found)
	if gotDir != wantDir {
		t.Errorf("expected %q, got %q", wantDir, gotDir)
	}
}

func TestFindLocalDir_NotFound(t *testing.T) {
	dir := t.TempDir()
	chdirTemp(t, dir)

	found, err := FindLocalDir()
	if err != nil {
		t.Fatal(err)
	}
	if found != "" {
		t.Errorf("expected empty string, got %q", found)
	}
}

// --- FindLocalConfigPath ---

func TestFindLocalConfigPath_Found(t *testing.T) {
	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	sub := filepath.Join(root, "nested", "subdir")
	os.MkdirAll(whoiamDir, 0700)
	os.MkdirAll(sub, 0700)
	os.WriteFile(filepath.Join(whoiamDir, "whoiam.yaml"), []byte("accounts: {}\n"), 0644)
	chdirTemp(t, sub)

	cp, err := FindLocalConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if cp == nil {
		t.Fatal("expected non-nil ConfigPath")
	}
	// Resolve symlinks (macOS /var/folders → /private/var/folders)
	wantPath, _ := filepath.EvalSymlinks(whoiamDir)
	gotPath, _ := filepath.EvalSymlinks(cp.Path)
	if gotPath != wantPath {
		t.Errorf("expected %q, got %q", wantPath, gotPath)
	}
}

func TestFindLocalConfigPath_NotFound(t *testing.T) {
	dir := t.TempDir()
	chdirTemp(t, dir)

	cp, err := FindLocalConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if cp != nil {
		t.Errorf("expected nil, got %+v", cp)
	}
}

// --- ReadCurrentEnvWithSource ---

func TestReadCurrentEnvWithSource_EnvVar(t *testing.T) {
	t.Setenv(ExpectedEnvVar, "production")

	name, source, err := ReadCurrentEnvWithSource()
	if err != nil {
		t.Fatal(err)
	}
	if name != "production" {
		t.Errorf("expected production, got %q", name)
	}
	if source != "env" {
		t.Errorf("expected env, got %q", source)
	}
}

func TestReadCurrentEnvWithSource_Local(t *testing.T) {
	t.Setenv(ExpectedEnvVar, "")

	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	os.WriteFile(filepath.Join(whoiamDir, "expected-env"), []byte("staging\n"), 0644)
	chdirTemp(t, root)

	name, source, err := ReadCurrentEnvWithSource()
	if err != nil {
		t.Fatal(err)
	}
	if name != "staging" {
		t.Errorf("expected staging, got %q", name)
	}
	if source != "local" {
		t.Errorf("expected local, got %q", source)
	}
}

func TestReadCurrentEnvWithSource_EmptyFile(t *testing.T) {
	t.Setenv(ExpectedEnvVar, "")

	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	// Write whitespace only — should be treated as not set
	os.WriteFile(filepath.Join(whoiamDir, "expected-env"), []byte("  \n"), 0644)
	chdirTemp(t, root)

	name, _, err := ReadCurrentEnvWithSource()
	if err != nil {
		t.Fatal(err)
	}
	// Falls through to global; we just verify it doesn't error and the local empty file is ignored
	_ = name
}

// --- WriteCurrentEnv / ReadCurrentEnv / ClearCurrentEnv ---

func TestWriteAndReadCurrentEnv(t *testing.T) {
	t.Setenv(ExpectedEnvVar, "")

	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	chdirTemp(t, root)

	if err := WriteCurrentEnv("dev"); err != nil {
		t.Fatalf("WriteCurrentEnv: %v", err)
	}

	name, err := ReadCurrentEnv()
	if err != nil {
		t.Fatal(err)
	}
	if name != "dev" {
		t.Errorf("expected dev, got %q", name)
	}
}

func TestWriteCurrentEnv_NoDir(t *testing.T) {
	dir := t.TempDir()
	chdirTemp(t, dir)

	if err := WriteCurrentEnv("dev"); err == nil {
		t.Error("expected error when no .whoiam/ dir exists")
	}
}

func TestClearCurrentEnv(t *testing.T) {
	t.Setenv(ExpectedEnvVar, "")

	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	os.WriteFile(filepath.Join(whoiamDir, "expected-env"), []byte("staging\n"), 0644)
	chdirTemp(t, root)

	if err := ClearCurrentEnv(); err != nil {
		t.Fatalf("ClearCurrentEnv: %v", err)
	}

	// File should be gone; ReadCurrentEnv should return "" (or fall to global — either way no local)
	_, err := os.Stat(filepath.Join(whoiamDir, "expected-env"))
	if !os.IsNotExist(err) {
		t.Error("expected expected-env file to be removed")
	}
}

func TestClearCurrentEnv_NoFile(t *testing.T) {
	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	chdirTemp(t, root)

	if err := ClearCurrentEnv(); err != nil {
		t.Fatalf("expected nil when file does not exist, got %v", err)
	}
}

func TestClearCurrentEnv_NoDir(t *testing.T) {
	dir := t.TempDir()
	chdirTemp(t, dir)

	// No .whoiam/ dir — should be a no-op, not an error
	if err := ClearCurrentEnv(); err != nil {
		t.Fatalf("expected nil when no .whoiam/ dir, got %v", err)
	}
}

// --- WriteGlobalCurrentEnv / ReadGlobalCurrentEnv / ClearGlobalCurrentEnv ---
// These use $HOME to locate ~/.whoiam. We redirect HOME to a temp dir so tests
// don't touch the real user config.

func TestReadGlobalCurrentEnv_Set(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	whoiamDir := filepath.Join(home, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	os.WriteFile(filepath.Join(whoiamDir, "expected-env"), []byte("prod\n"), 0644)

	name, err := ReadGlobalCurrentEnv()
	if err != nil {
		t.Fatal(err)
	}
	if name != "prod" {
		t.Errorf("expected prod, got %q", name)
	}
}

func TestReadGlobalCurrentEnv_NotSet(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	os.MkdirAll(filepath.Join(home, ".whoiam"), 0700)

	name, err := ReadGlobalCurrentEnv()
	if err != nil {
		t.Fatal(err)
	}
	if name != "" {
		t.Errorf("expected empty string, got %q", name)
	}
}

func TestWriteAndReadGlobalCurrentEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	os.MkdirAll(filepath.Join(home, ".whoiam"), 0700)

	if err := WriteGlobalCurrentEnv("prod"); err != nil {
		t.Fatalf("WriteGlobalCurrentEnv: %v", err)
	}

	name, err := ReadGlobalCurrentEnv()
	if err != nil {
		t.Fatal(err)
	}
	if name != "prod" {
		t.Errorf("expected prod, got %q", name)
	}
}

func TestWriteGlobalCurrentEnv_NoDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	// No .whoiam/ dir created

	if err := WriteGlobalCurrentEnv("prod"); err == nil {
		t.Error("expected error when ~/.whoiam/ does not exist")
	}
}

func TestClearGlobalCurrentEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	whoiamDir := filepath.Join(home, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	os.WriteFile(filepath.Join(whoiamDir, "expected-env"), []byte("prod\n"), 0644)

	if err := ClearGlobalCurrentEnv(); err != nil {
		t.Fatalf("ClearGlobalCurrentEnv: %v", err)
	}

	_, err := os.Stat(filepath.Join(whoiamDir, "expected-env"))
	if !os.IsNotExist(err) {
		t.Error("expected expected-env to be removed")
	}
}

func TestClearGlobalCurrentEnv_NoFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	os.MkdirAll(filepath.Join(home, ".whoiam"), 0700)

	if err := ClearGlobalCurrentEnv(); err != nil {
		t.Fatalf("expected nil when file does not exist, got %v", err)
	}
}

// --- LoadEffectiveConfig ---

func TestLoadEffectiveConfig_LocalOnly(t *testing.T) {
	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	cp := &ConfigPath{Path: whoiamDir, File: "whoiam.yaml"}
	cp.SaveConfig(&Config{Accounts: map[string]string{"dev": "111111111111"}})
	chdirTemp(t, root)

	loaded, err := LoadEffectiveConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.AccountExists("dev") {
		t.Error("expected 'dev' account to exist")
	}
	if loaded.Accounts["dev"] != "111111111111" {
		t.Errorf("expected 111111111111, got %q", loaded.Accounts["dev"])
	}
}

func TestLoadEffectiveConfigWithSources_LocalSourceLabeled(t *testing.T) {
	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	cp := &ConfigPath{Path: whoiamDir, File: "whoiam.yaml"}
	cp.SaveConfig(&Config{Accounts: map[string]string{
		"dev":     "222222222222",
		"staging": "333333333333",
	}})
	chdirTemp(t, root)

	_, sources, err := LoadEffectiveConfigWithSources()
	if err != nil {
		t.Fatal(err)
	}
	if sources["dev"] != "local" {
		t.Errorf("expected source 'local' for dev, got %q", sources["dev"])
	}
	if sources["staging"] != "local" {
		t.Errorf("expected source 'local' for staging, got %q", sources["staging"])
	}
}

func TestLoadEffectiveConfigWithSources_GlobalAndLocal(t *testing.T) {
	// Global config in fake HOME
	home := t.TempDir()
	t.Setenv("HOME", home)
	globalDir := filepath.Join(home, ".whoiam")
	os.MkdirAll(globalDir, 0700)
	globalCp := &ConfigPath{Path: globalDir, File: "whoiam.yaml"}
	globalCp.SaveConfig(&Config{Accounts: map[string]string{
		"prod": "111111111111",
		"dev":  "222222222222", // will be overridden by local
	}})

	// Local config in temp project dir
	root := t.TempDir()
	localDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(localDir, 0700)
	localCp := &ConfigPath{Path: localDir, File: "whoiam.yaml"}
	localCp.SaveConfig(&Config{Accounts: map[string]string{
		"dev": "999999999999", // overrides global dev
	}})
	chdirTemp(t, root)

	cfg, sources, err := LoadEffectiveConfigWithSources()
	if err != nil {
		t.Fatal(err)
	}
	if sources["prod"] != "global" {
		t.Errorf("expected source 'global' for prod, got %q", sources["prod"])
	}
	if sources["dev"] != "local" {
		t.Errorf("expected source 'local' for dev (local override), got %q", sources["dev"])
	}
	if cfg.Accounts["dev"] != "999999999999" {
		t.Errorf("expected local dev account to override global, got %q", cfg.Accounts["dev"])
	}
}

func TestReadCurrentEnvWithSource_GlobalFallback(t *testing.T) {
	t.Setenv(ExpectedEnvVar, "")

	home := t.TempDir()
	t.Setenv("HOME", home)
	whoiamDir := filepath.Join(home, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	os.WriteFile(filepath.Join(whoiamDir, "expected-env"), []byte("global-env\n"), 0644)

	// CWD has no local .whoiam/
	chdirTemp(t, t.TempDir())

	name, source, err := ReadCurrentEnvWithSource()
	if err != nil {
		t.Fatal(err)
	}
	if name != "global-env" {
		t.Errorf("expected global-env, got %q", name)
	}
	if source != "global" {
		t.Errorf("expected source global, got %q", source)
	}
}

// --- LoadConfig invalid YAML ---

func TestConfigPathLoadConfig_InvalidYAML(t *testing.T) {
	cp := newTempConfigPath(t)
	cp.Create()
	os.WriteFile(cp.FullPath(), []byte("accounts: {invalid yaml:: ::!!}"), 0644)

	_, err := cp.LoadConfig()
	if err == nil {
		t.Error("expected error for invalid YAML content")
	}
}

// --- PrintConfigTable / PrintConfigTableWithSource ---


func TestPrintConfigTable(t *testing.T) {
	cfg := &Config{Accounts: map[string]string{"dev": "111111111111"}}
	var buf bytes.Buffer
	cfg.PrintConfigTable(&buf)
	out := buf.String()
	if !strings.Contains(out, "dev") {
		t.Errorf("expected output to contain account name, got: %s", out)
	}
	if !strings.Contains(out, "111111111111") {
		t.Errorf("expected output to contain account number, got: %s", out)
	}
}

func TestPrintConfigTableWithSource_WithSources(t *testing.T) {
	cfg := &Config{Accounts: map[string]string{"prod": "222222222222"}}
	sources := map[string]string{"prod": "global"}
	var buf bytes.Buffer
	cfg.PrintConfigTableWithSource(&buf, sources)
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Errorf("expected output to contain account name, got: %s", out)
	}
	if !strings.Contains(out, "global") {
		t.Errorf("expected output to contain source, got: %s", out)
	}
}

func TestPrintConfigTableWithSource_NilSources(t *testing.T) {
	cfg := &Config{Accounts: map[string]string{"staging": "333333333333"}}
	var buf bytes.Buffer
	cfg.PrintConfigTableWithSource(&buf, nil)
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected output to contain account name, got: %s", out)
	}
}
