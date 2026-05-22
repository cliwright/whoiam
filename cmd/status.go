/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"github.com/pytoolbelt/whoiam/internal"
	"github.com/spf13/cobra"
)

func statusEntrypoint(cmd *cobra.Command, args []string) error {
	currentEnv, source, err := internal.ReadCurrentEnvWithSource()
	if err != nil {
		return err
	}

	if currentEnv == "" {
		cmd.Println("Expected env: not set")
	} else {
		cmd.Printf("Expected env: %s (%s)\n", currentEnv, source)
	}

	client, err := internal.NewStsClient()
	if err != nil {
		cmd.Println("Authenticated:  no — could not create AWS client")
		return nil
	}

	identity, err := client.GetCallerIdentity()
	if err != nil {
		cmd.Println("Authenticated:  no — not authenticated")
		return nil
	}

	cfg, err := internal.LoadEffectiveConfig()
	if err != nil {
		return err
	}

	accountName := cfg.GetAccountByNumber(*identity.Account)
	if accountName == "" {
		accountName = "unknown"
	}

	cmd.Printf("Authenticated:  yes\n")
	cmd.Printf("Account:        %s (%s)\n", accountName, *identity.Account)
	cmd.Printf("ARN:            %s\n", *identity.Arn)
	return nil
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current environment and authenticated AWS account",
	RunE:  statusEntrypoint,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
