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

func execEntrypoint(cmd *cobra.Command, args []string) {
	accountName, _ := cmd.Flags().GetString("account")

	if accountName == "" {
		fmt.Println("Account Name is required")
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Println("No command provided starting subshell. Type 'exit' to return to the parent shell")
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

	err = internal.AssertAccountAsExpected(identity, cfg.Accounts[accountName])
	internal.HandleError(err)
	fmt.Printf("Verified: AWS account %s (%s)\n", accountName, cfg.Accounts[accountName])

	shell, err := internal.NewSubShell(args...)
	internal.HandleError(err)

	err = shell.Run()
	internal.HandleError(err)
}

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Verify the current AWS account and run a command (or open a subshell)",
	Long: `Verifies that the current AWS credentials match the expected account before
running the given command. If no command is provided, opens an interactive
subshell with the account assertion already satisfied.

Example:
  whoiam exec --account production terraform apply
  whoiam exec --account staging`,
	Run: execEntrypoint,
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringP("account", "a", "", "Account Name")
}
