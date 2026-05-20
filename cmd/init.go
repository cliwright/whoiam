/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/whoiam/internal"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func projectInitEntrypoint(cmd *cobra.Command, args []string) {
	global, _ := cmd.Flags().GetBool("global")

	if global {
		globalPath, err := internal.NewConfigPath()
		internal.HandleError(err)

		if globalPath.ConfigFileExists() {
			fmt.Printf("Global config already exists at %s\n", globalPath.FullPath())
			os.Exit(1)
		}

		err = globalPath.Create()
		internal.HandleError(err)

		cfg, err := internal.NewTemplateConfig()
		internal.HandleError(err)

		err = globalPath.SaveConfig(cfg)
		internal.HandleError(err)

		fmt.Printf("Initialized global config at %s\n", globalPath.FullPath())
		return
	}

	localPath, err := internal.NewProjectConfigPath()
	internal.HandleError(err)

	if localPath.Exists() {
		fmt.Printf("Project config already exists at %s\n", localPath.FullPath())
		os.Exit(1)
	}

	err = localPath.Create()
	internal.HandleError(err)

	cfg, err := internal.NewTemplateConfig()
	internal.HandleError(err)

	err = localPath.SaveConfig(cfg)
	internal.HandleError(err)

	// Create .gitignore to exclude expected-env (session state, not shared config)
	gitignorePath := filepath.Join(localPath.Path, ".gitignore")
	err = os.WriteFile(gitignorePath, []byte("expected-env\n"), 0644)
	internal.HandleError(err)

	fmt.Printf("Initialized project config at %s\n", localPath.FullPath())
	fmt.Println("Edit the config file to add your project's AWS account mappings.")
	fmt.Println("Tip: commit .whoiam/whoiam.yaml to share account mappings with your team.")
}

var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a whoiam config (project-local by default, --global for ~/.whoiam)",
	Long: `Creates a .whoiam/ directory in the current directory with a template config
file and a .gitignore that excludes session state (expected-env).

Pass --global to initialize the global config at ~/.whoiam/whoiam.yaml instead.

Commit .whoiam/whoiam.yaml to share account name mappings with your team.
Do not commit .whoiam/expected-env — it is personal session state.`,
	Run: projectInitEntrypoint,
}

func init() {
	rootCmd.AddCommand(projectInitCmd)
	projectInitCmd.Flags().Bool("global", false, "Initialize the global config at ~/.whoiam/whoiam.yaml")
}
