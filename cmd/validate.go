/*
Copyright © 2024 Jesse Maitland jesse@cliwright.com
*/
package cmd

import (
	"fmt"
	"github.com/cliwright/whoiam/internal"
	"github.com/spf13/cobra"
)

func validateEntrypoint(cmd *cobra.Command, args []string) error {
	accountName, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	cfg, err := internal.LoadEffectiveConfig()
	if err != nil {
		return err
	}

	if accountName == "" {
		envAccount, _, err := internal.ReadCurrentEnvForConfig(cfg)
		if err != nil {
			return err
		}
		accountName = envAccount
	}

	if accountName == "" {
		return fmt.Errorf("no account specified — use --env or run 'whoiam set <account>'")
	}

	if !cfg.AccountExists(accountName) {
		return fmt.Errorf("account %q does not exist in config", accountName)
	}

	client, err := internal.NewStsClient()
	if err != nil {
		return err
	}

	identity, err := client.GetCallerIdentity()
	if err != nil {
		return err
	}

	if err := internal.AssertAccountAsExpected(identity, cfg.Accounts[accountName]); err != nil {
		return fmt.Errorf("account mismatch: expected %s (%s), got %s", accountName, cfg.Accounts[accountName], *identity.Account)
	}

	cmd.Printf("OK: %s (%s)\n", accountName, cfg.Accounts[accountName])
	return nil
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
	RunE: validateEntrypoint,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("env", "e", "", "Expected environment name")
}
