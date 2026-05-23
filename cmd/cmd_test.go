package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cliwright/whoiam/internal"
)

// executeCmd runs a cobra command with the given args and returns stdout output and any error.
func executeCmd(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

// setupProjectDir creates a temp dir with a .whoiam/whoiam.yaml containing the given accounts,
// changes CWD to it, and restores CWD when the test ends.
func setupProjectDir(t *testing.T, accounts map[string]string) string {
	t.Helper()
	root := t.TempDir()
	whoiamDir := filepath.Join(root, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	cp := &internal.ConfigPath{Path: whoiamDir, File: "whoiam.yaml"}
	cp.SaveConfig(&internal.Config{Accounts: accounts})

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
	return root
}

// --- whoiam init ---

func TestInitCmd_Local(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(orig) })

	out, err := executeCmd("init")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Initialized project config") {
		t.Errorf("expected success message, got: %s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, ".whoiam", "whoiam.yaml")); os.IsNotExist(err) {
		t.Error("expected whoiam.yaml to be created")
	}
	if _, err := os.Stat(filepath.Join(dir, ".whoiam", ".gitignore")); os.IsNotExist(err) {
		t.Error("expected .gitignore to be created")
	}
}

func TestInitCmd_Local_AlreadyExists(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(orig) })

	// First init
	executeCmd("init")

	// Second init should fail
	_, err := executeCmd("init")
	if err == nil {
		t.Error("expected error when project config already exists")
	}
}

func TestInitCmd_Global(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	out, err := executeCmd("init", "--global")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Initialized global config") {
		t.Errorf("expected success message, got: %s", out)
	}
	if _, err := os.Stat(filepath.Join(home, ".whoiam", "whoiam.yaml")); os.IsNotExist(err) {
		t.Error("expected global whoiam.yaml to be created")
	}
}

func TestInitCmd_Global_AlreadyExists(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	executeCmd("init", "--global")

	_, err := executeCmd("init", "--global")
	if err == nil {
		t.Error("expected error when global config already exists")
	}
}

// --- whoiam set ---

func TestSetCmd_SetLocal(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	setupProjectDir(t, map[string]string{"dev": "111111111111"})

	out, err := executeCmd("set", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "dev") {
		t.Errorf("expected confirmation message, got: %s", out)
	}

	name, err := internal.ReadCurrentEnv()
	if err != nil {
		t.Fatal(err)
	}
	if name != "dev" {
		t.Errorf("expected dev, got %q", name)
	}
}

func TestSetCmd_SetLocal_UnknownAccount(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	setupProjectDir(t, map[string]string{"dev": "111111111111"})

	_, err := executeCmd("set", "nonexistent")
	if err == nil {
		t.Error("expected error for unknown account")
	}
}

func TestSetCmd_ClearLocal(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	root := setupProjectDir(t, map[string]string{"dev": "111111111111"})
	// Pre-write an expected-env file
	os.WriteFile(filepath.Join(root, ".whoiam", "expected-env"), []byte("dev\n"), 0644)

	out, err := executeCmd("set")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Cleared") {
		t.Errorf("expected cleared message, got: %s", out)
	}
}

func TestSetCmd_SetGlobal(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	home := t.TempDir()
	t.Setenv("HOME", home)
	os.MkdirAll(filepath.Join(home, ".whoiam"), 0700)

	// Write global config with account
	globalCp := &internal.ConfigPath{Path: filepath.Join(home, ".whoiam"), File: "whoiam.yaml"}
	globalCp.SaveConfig(&internal.Config{Accounts: map[string]string{"prod": "222222222222"}})

	out, err := executeCmd("set", "--global", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("expected confirmation message, got: %s", out)
	}
}

func TestSetCmd_ClearGlobal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	whoiamDir := filepath.Join(home, ".whoiam")
	os.MkdirAll(whoiamDir, 0700)
	os.WriteFile(filepath.Join(whoiamDir, "expected-env"), []byte("prod\n"), 0644)

	out, err := executeCmd("set", "--global")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "Cleared") {
		t.Errorf("expected cleared message, got: %s", out)
	}
}

// --- whoiam config ---

func TestConfigCmd(t *testing.T) {
	setupProjectDir(t, map[string]string{"dev": "111111111111", "prod": "222222222222"})

	out, err := executeCmd("config")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "dev") || !strings.Contains(out, "prod") {
		t.Errorf("expected account names in output, got: %s", out)
	}
}

// --- whoiam status ---

func TestStatusCmd_NoEnvSet(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(orig) })

	out, err := executeCmd("status")
	// status is graceful — never returns an error even when not authenticated
	if err != nil {
		t.Fatalf("status should not error, got: %v", err)
	}
	if !strings.Contains(out, "Expected env: not set") {
		t.Errorf("expected 'not set' in output, got: %s", out)
	}
}

func TestStatusCmd_EnvVarSet(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "staging")
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(orig) })

	out, err := executeCmd("status")
	if err != nil {
		t.Fatalf("status should not error, got: %v", err)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected staging in output, got: %s", out)
	}
}

// --- whoiam exec ---

func TestExecCmd_NoAccountSpecified(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(orig) })

	_, err := executeCmd("exec")
	if err == nil {
		t.Error("expected error when no account specified")
	}
}

func TestExecCmd_UnknownAccount(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	setupProjectDir(t, map[string]string{"dev": "111111111111"})

	_, err := executeCmd("exec", "--env", "nonexistent")
	if err == nil {
		t.Error("expected error for unknown account")
	}
}

func TestExecCmd_NestedSubshell(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	t.Setenv(internal.SubShellVar, "1")
	setupProjectDir(t, map[string]string{"dev": "111111111111"})

	_, err := executeCmd("exec", "--env", "dev")
	if err == nil {
		t.Error("expected error when already in a subshell")
	}
}

// --- whoiam validate ---

func TestValidateCmd_NoAccountSpecified(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(orig) })

	_, err := executeCmd("validate")
	if err == nil {
		t.Error("expected error when no account specified")
	}
}

func TestValidateCmd_UnknownAccount(t *testing.T) {
	t.Setenv(internal.ExpectedEnvVar, "")
	setupProjectDir(t, map[string]string{"dev": "111111111111"})

	_, err := executeCmd("validate", "--env", "nonexistent")
	if err == nil {
		t.Error("expected error for unknown account")
	}
}
