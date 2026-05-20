/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/whoiam/internal"
	"github.com/spf13/cobra"
)

func statusEntrypoint(cmd *cobra.Command, args []string) {
	currentEnv, source, err := internal.ReadCurrentEnvWithSource()
	internal.HandleError(err)

	if currentEnv == "" {
		fmt.Println("Expected env: not set")
	} else {
		fmt.Printf("Expected env: %s (%s)\n", currentEnv, source)
	}

	client, err := internal.NewStsClient()
	if err != nil {
		fmt.Println("Authenticated:  no — could not create AWS client")
		return
	}

	identity, err := client.GetCallerIdentity()
	if err != nil {
		fmt.Println("Authenticated:  no — not authenticated")
		return
	}

	cfg, err := internal.LoadEffectiveConfig()
	internal.HandleError(err)

	accountName := cfg.GetAccountByNumber(*identity.Account)
	if accountName == "" {
		accountName = "unknown"
	}

	fmt.Printf("Authenticated:  yes\n")
	fmt.Printf("Account:        %s (%s)\n", accountName, *identity.Account)
	fmt.Printf("ARN:            %s\n", *identity.Arn)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current environment and authenticated AWS account",
	Run:   statusEntrypoint,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
