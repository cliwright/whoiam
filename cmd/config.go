/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"github.com/pytoolbelt/whoiam/internal"
	"github.com/spf13/cobra"
	"os"
)

func configEntrypoint(cmd *cobra.Command, args []string) {
	cfg, sources, err := internal.LoadEffectiveConfigWithSources()
	internal.HandleError(err)

	cfg.PrintConfigTableWithSource(sources)

	os.Exit(0)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the effective whoiam config (global + project-local merged)",
	Run:   configEntrypoint,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
