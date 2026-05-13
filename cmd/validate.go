/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"
	"github.com/pytoolbelt/whoiam/internal"
	"github.com/spf13/cobra"
	"os"
)

func validateEntrypoint(cmd *cobra.Command, args []string) {
	accountName, _ := cmd.Flags().GetString("env")

	if accountName == "" {
		envAccount, err := internal.ReadCurrentEnv()
		internal.HandleError(err)
		accountName = envAccount
	}

	if accountName == "" {
		fmt.Println("No account specified — use --env or run 'whoiam set <account>'")
		os.Exit(1)
	}

	cfg, err := internal.LoadEffectiveConfig()
	internal.HandleError(err)

	if !cfg.AccountExists(accountName) {
		fmt.Printf("Account %q does not exist in config\n", accountName)
		os.Exit(1)
	}

	client, err := internal.NewStsClient()
	internal.HandleError(err)

	identity, err := client.GetCallerIdentity()
	internal.HandleError(err)

	expectedNumber := cfg.Accounts[accountName]
	if *identity.Account != expectedNumber {
		fmt.Printf("Account mismatch: expected %s (%s), got %s\n", accountName, expectedNumber, *identity.Account)
		os.Exit(1)
	}

	fmt.Printf("OK: %s (%s)\n", accountName, expectedNumber)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate that the current AWS account matches the expected environment",
	Long: `Compares the current AWS caller identity against the expected account.
Exits with a non-zero status if they do not match, making it safe to use
as a pre-flight check in Taskfile, mise, or CI pipelines.

The expected account is resolved from --env if provided, otherwise from
the project session set by 'whoiam set'.

Examples:
  whoiam validate                    # uses current-env
  whoiam validate --env production   # explicit environment`,
	Run: validateEntrypoint,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("env", "e", "", "Expected environment name")
}
