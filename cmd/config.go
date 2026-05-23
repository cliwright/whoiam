/*
Copyright © 2024 Jesse Maitland jesse@cliwright.com
*/
package cmd

import (
	"github.com/cliwright/whoiam/internal"
	"github.com/spf13/cobra"
)

func configEntrypoint(cmd *cobra.Command, args []string) error {
	cfg, sources, err := internal.LoadEffectiveConfigWithSources()
	if err != nil {
		return err
	}
	cfg.PrintConfigTableWithSource(cmd.OutOrStdout(), sources)
	return nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the effective whoiam config (global + project-local merged)",
	RunE:  configEntrypoint,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
