/*
Copyright © 2024 Jesse Maitland jesse@cliwright.com
*/
package cmd

import (
	"github.com/cliwright/whoiam/internal"
	"os"

	"github.com/spf13/cobra"
)

func rootEntrypoint(cmd *cobra.Command, args []string) error {
	client, err := internal.NewStsClient()
	if err != nil {
		return err
	}

	identity, err := client.GetCallerIdentity()
	if err != nil {
		return err
	}

	cfg, err := internal.LoadEffectiveConfig()
	if err != nil {
		return err
	}

	accountName := cfg.GetAccountByNumber(*identity.Account)
	if accountName == "" {
		accountName = "Unknown"
	}

	internal.PrintCallerIdentityTable(identity, accountName)
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "whoiam",
	Short:        "Check your current AWS IAM Role",
	Long:         ``,
	RunE:         rootEntrypoint,
	SilenceUsage: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
