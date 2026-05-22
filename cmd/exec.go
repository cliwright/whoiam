/*
Copyright © 2024 Jesse Maitland jesse@pytoolbelt.com
*/
package cmd

import (
	"fmt"

	"github.com/pytoolbelt/whoiam/internal"
	"github.com/spf13/cobra"
)

func execEntrypoint(cmd *cobra.Command, args []string) error {
	accountName, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	if accountName == "" {
		envAccount, err := internal.ReadCurrentEnv()
		if err != nil {
			return err
		}
		accountName = envAccount
	}

	if accountName == "" {
		return fmt.Errorf("no account specified — use --env or run 'whoiam set <account>'")
	}

	if len(args) == 0 {
		cmd.Println("No command provided starting subshell. Type 'exit' to return to the parent shell")
	}

	cfg, err := internal.LoadEffectiveConfig()
	if err != nil {
		return err
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
		return err
	}
	cmd.Printf("Verified: AWS account %s (%s)\n", accountName, cfg.Accounts[accountName])

	shell, err := internal.NewSubShell(args...)
	if err != nil {
		return err
	}

	return shell.Run()
}

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Verify the current AWS account and run a command (or open a subshell)",
	Long: `Verifies that the current AWS credentials match the expected environment before
running the given command. If no command is provided, opens an interactive
subshell with the account assertion already satisfied.

Example:
  whoiam exec --env production terraform apply
  whoiam exec --env staging`,
	RunE: execEntrypoint,
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringP("env", "e", "", "Environment name")
}
