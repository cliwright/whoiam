/*
Copyright © 2024 Jesse Maitland jesse@cliwright.com
*/
package cmd

import (
	"fmt"
	"github.com/cliwright/whoiam/internal"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func projectInitEntrypoint(cmd *cobra.Command, args []string) error {
	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	if global {
		globalPath, err := internal.NewConfigPath()
		if err != nil {
			return err
		}

		if globalPath.ConfigFileExists() {
			return fmt.Errorf("global config already exists at %s", globalPath.FullPath())
		}

		if err := globalPath.Create(); err != nil {
			return err
		}

		cfg, err := internal.NewTemplateConfig()
		if err != nil {
			return err
		}

		if err := globalPath.SaveConfig(cfg); err != nil {
			return err
		}

		cmd.Printf("Initialized global config at %s\n", globalPath.FullPath())
		return nil
	}

	localPath, err := internal.NewProjectConfigPath()
	if err != nil {
		return err
	}

	if localPath.Exists() {
		return fmt.Errorf("project config already exists at %s", localPath.FullPath())
	}

	if err := localPath.Create(); err != nil {
		return err
	}

	cfg, err := internal.NewTemplateConfig()
	if err != nil {
		return err
	}

	if err := localPath.SaveConfig(cfg); err != nil {
		return err
	}

	// Create .gitignore to exclude expected-env (session state, not shared config)
	gitignorePath := filepath.Join(localPath.Path, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("expected-env\n"), 0644); err != nil {
		return err
	}

	cmd.Printf("Initialized project config at %s\n", localPath.FullPath())
	cmd.Println("Edit the config file to add your project's AWS account mappings.")
	cmd.Println("Tip: commit .whoiam/whoiam.yaml to share account mappings with your team.")
	return nil
}

var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a whoiam config (project-local by default, --global for ~/.whoiam)",
	Long: `Creates a .whoiam/ directory in the current directory with a template config
file and a .gitignore that excludes session state (expected-env).

Pass --global to initialize the global config at ~/.whoiam/whoiam.yaml instead.

Commit .whoiam/whoiam.yaml to share account name mappings with your team.
Do not commit .whoiam/expected-env — it is personal session state.`,
	RunE: projectInitEntrypoint,
}

func init() {
	rootCmd.AddCommand(projectInitCmd)
	projectInitCmd.Flags().Bool("global", false, "Initialize the global config at ~/.whoiam/whoiam.yaml")
}
